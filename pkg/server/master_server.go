package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/yiqinguo/armyant/pkg/models"
	"github.com/yiqinguo/armyant/pkg/server/storage"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
)

const (
	namespace string = "armyant"
)

type benchmarkServer struct {
	storage storage.Storage
	kClient *kubernetes.Clientset
}

func (b *benchmarkServer) ReportResult(ctx context.Context,
	stats *models.Stats) (*Response, error) {

	resp := &Response{Code: 1}

	body, _ := json.Marshal(stats)
	log.Printf("%s", body)

	err := b.storage.AddResult(stats)
	if err != nil {
		resp.Message = err.Error()
		return resp, err
	}
	resp.Code = 0
	return resp, nil
}

func (b *benchmarkServer) ReportStatus(ctx context.Context,
	stats *models.Status) (*Response, error) {
	log.Println("=============")

	return &Response{Code: 0}, nil
}

func (b *benchmarkServer) CreateJob(ctx context.Context,
	jobSpec *models.Job) (*Response, error) {

	resp := &Response{
		Code: 1,
	}
	job, err := b.kClient.BatchV1().Jobs(namespace).Get("testjob", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		resp.Message = err.Error()
		return resp, err
	}

	var replica int32 = getReplicas(jobSpec.Connections)
	requests := jobSpec.Requests / int64(replica)
	connections := jobSpec.Connections / int64(replica)

	job = &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("armyant-%s", jobSpec.Id),
			Namespace: namespace,
		},
		Spec: batchv1.JobSpec{
			Parallelism: &replica,
			Completions: &replica,
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name:    "job",
							Image:   "ruly-reg.ofo.com/ruly/armyant-slave:latest",
							Command: []string{"./slave-agent"},
							Args: []string{
								fmt.Sprintf("--id=%s", jobSpec.Id),
								fmt.Sprintf("--hid=%s", jobSpec.HistoryId),
								fmt.Sprintf("--url=%s", jobSpec.Url),
								fmt.Sprintf("--requests=%d", requests),
								fmt.Sprintf("--master=%s", jobSpec.Master),
								fmt.Sprintf("--conn=%d", connections),
							},
							ImagePullPolicy: v1.PullAlways,
							Env: []v1.EnvVar{
								v1.EnvVar{
									Name: "POD_NAME",
									ValueFrom: &v1.EnvVarSource{
										FieldRef: &v1.ObjectFieldSelector{
											FieldPath: "metadata.name",
										},
									},
								},
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
					ImagePullSecrets: []v1.LocalObjectReference{
						v1.LocalObjectReference{Name: "ruly-reg.ofo.com"},
					},
				},
			},
		},
	}
	_, err = b.kClient.BatchV1().Jobs(namespace).Create(job)
	if err != nil {
		return resp, err
	}

	return &Response{Code: 0}, nil
}

func (b *benchmarkServer) UpdateJob(ctx context.Context,
	stats *models.Job) (*Response, error) {
	log.Println("=============")

	return &Response{Code: 200}, nil
}

func (b *benchmarkServer) ListJobs(ctx context.Context,
	listJobRequest *models.ListJobRequest) (*Response, error) {
	log.Println("=============")

	return &Response{Code: 200}, nil
}

func (b *benchmarkServer) GetJob(ctx context.Context,
	getJobRequest *models.GetJobRequest) (*models.JobResultResponse, error) {

	response := &models.JobResultResponse{
		Code:   1,
		Result: map[string]float32{},
	}
	/*if getJobRequest.JobId == "" {
		return response, fmt.Errorf("job id cannot be empty")
	}

	resp, err := b.storage.GetJobResult(getJobRequest.JobId, getJobRequest.InstancdId)
	if err != nil {
		return response, err
	}
	response.Code = int64(0)
	response.Result = resp.Result*/

	return response, nil
}

func (b *benchmarkServer) GetJobResult(ctx context.Context,
	getJobRequest *models.GetJobRequest) (*models.JobResultResponse, error) {

	response := &models.JobResultResponse{
		Code:   1,
		Result: map[string]float32{},
	}
	if getJobRequest.JobId == "" {
		return response, fmt.Errorf("job id cannot be empty")
	}

	resp, err := b.storage.GetJobResult(getJobRequest.JobId, getJobRequest.InstancdId)
	if err != nil {
		return response, err
	}
	response.Code = int64(0)
	response.Result = resp.Result

	return response, nil
}

func (b *benchmarkServer) Run(args models.MasterConfig) error {
	log.Println("grpc server starting...")
	listen, err := net.Listen("tcp", ":"+args.GrpcPort)
	if err != nil {
		log.Printf("listen tcp port error: %v", err)
		return err
	}

	s := grpc.NewServer()
	RegisterApiserverServer(s, b)
	reflection.Register(s)

	err = s.Serve(listen)
	if err != nil {
		log.Printf("start grpc server error: %v", err)
	}
	return err
}

func NewBenchmarkServer(args models.MasterConfig) *benchmarkServer {
	log.Println(args.Storage)
	bs := &benchmarkServer{}
	switch args.Storage {
	case "local":
		/*var err error
		bs.storage, err = storage.NewBboltStorage(args.DataDir)
		if err != nil {
			log.Fatalf("new bolt db error: %v", err)
		}*/
	case "influxdb":
		var err error
		bs.storage, err = storage.NewInfluxdbStorage(
			args.InfluxUrl,
			args.InfluxDB,
			args.InfluxUsername,
			args.InfluxPassword,
		)
		if err != nil {
			log.Fatalf("new bolt db error: %v", err)
		}
	default:
		log.Fatal("invalid storage type")
	}

	kClient, err := newKubernetsClient()
	if err != nil {
		log.Fatalf("new kubernetes client error: %v", err)
	}
	bs.kClient = kClient

	return bs
}

func newKubernetsClient() (*kubernetes.Clientset, error) {
	//config, err := clientcmd.BuildConfigFromFlags("", "/Users/yiqinguo/.kube/config")
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getReplicas(conns int64) int32 {
	ms := conns / 100
	if ms >= 10 {
		return int32(10)
	} else if ms >= 1 {
		return int32(ms)
	}

	return int32(1)
}
