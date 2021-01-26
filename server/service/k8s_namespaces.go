package service

import (
	"context"
	"gin-vue-admin/global"
	"gin-vue-admin/model"
	"gin-vue-admin/model/request"
	"gin-vue-admin/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

// 定义结构体绑定数据
type Result struct {
	Namespace string            `json:"namespace"`
	Status    v1.NamespacePhase `json:"status"`
	Time      metav1.Time       `json:"time"`
}

//@function: GetK8sNamespaces
//@description: 根据id获取K8sNamespaces记录
//@param: id uint
//@return: err error, k8sNamespaces model.K8sNamespaces
func FindK8sNamespaces(id uint) (err error, k8sNamespaces model.K8sNamespaces) {
	err = global.GVA_DB.Where("id = ?", id).First(&k8sNamespaces).Error
	return
}

//@function: GetK8sNamespacesList
//@description: 分页获取K8sNamespaces记录
//@param: info request.K8sNamespacesSearch
//@return: err error, list []*Result, total int64
func GetK8sNamespacesList(info request.K8sNamespacesSearch) (err error, list []*Result, total int64) {
	// 初始化k8s客户端
	clientset, err := utils.InitClient()
	if err != nil {
		log.Fatalln(err)
	}
	// 获取所有Namespaces
	ns, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for _, nss := range ns.Items {
		res := &Result{
			Namespace: nss.Name,
			Status:    nss.Status.Phase,
			Time:      nss.CreationTimestamp,
		}
		list = append(list, res)
	}
	total = int64(len(list))
	return err, list, total
}