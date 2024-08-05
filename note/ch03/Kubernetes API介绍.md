# 一. 认识Kubernetes API

我们知道Kubernetes集群中有一个kube-apiserver组件，所有组件需要操作集群资源时都通过调用kube-apiserver提供的RESTful接口来实现。kube-apiserver进一步和ETCD交互，完成资源信息的更新，另外kube-apiserver也是集群内唯一和ETCD直接交互的组件。

Kubernetes中的资源本质就是一个API对象，这个对象的“期望状态”被API Server保存在ETCD中，然后提供RESTful接口用于更新这些对象。我们可以直接和API Server交互，使用“声明”的方式来管理这些资源（API对象）​，也可以通过kubectl这种命令行工具，或者client-go这类SDK。当然，在Dashboard上操作资源的方式和API Server交互也是一种方式，不过这不是我们关注的重点。

我们可以在Kubernetes的官网中找到API文档：[Kubernetes API文档](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/)。

# 二. Curl方式访问API

我们以一个Deployment资源的增删改查操作为例，看一下如何与API Server交互。

## 1.准备工作

由于kube-apiserver默认提供的是HTTPS服务，而且是双向TLS认证，而我们目前的关注重点是API本身，因此先通过kubectl来代理API Server服务

```shell
#可以通过简单的HTTP请求来和API Server交互
kubectl proxy --port=8080
#测试
curl localhost:8080/version
```

需要一个配置文件来描述Deployment资源，在本地创建一个nginx-deploy.yaml文件

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deploy
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.14
          ports:
            - containerPort: 80
```

## 2.资源创建

```shell
#在default命名空间下创建一个Deployment
curl -X POST \
-H 'Content-Type: application/yaml' \
--data-binary '@nginx-deploy.yaml' \
http://localhost:8080/apis/apps/v1/namespaces/default/deployments
#通过kubectl来查询这个资源对象
kubectl get deployment
```

## 3.资源删除

```shell
#在default命名空间下创建一个Deployment
curl -X DELETE -H 'Content-Type: application/yaml' \
--data-binary '
gracePeriodSeconds: 0
orphanDependents: false
' \
http://localhost:8080/apis/apps/v1/namespaces/default/deployments/nginx-deploy
#通过kubectl来查询这个资源对象
kubectl get deployment
```

# 三. kubectl raw方式访问API

通过kubectl get --raw可以实现和curl类似的效果，只是我们不需要指定API Server的地址了，同样认证信息也不需要，这时默认用了kubeconfig中的连接信息。

```shell
#访问API
kubectl get --raw /version
#查询nginx-deployment：
kubectl get --raw apis/apps/v1/namespaces/default/deployments/nginx-deploy
```

# 四. 理解GVK：组、版本与类型

其实我们前面已经多次接触到GVK这个词了，望文生义，大家不难猜到GVK就是<mark>Group</mark>、<mark>Version</mark>、<mark>Kind</mark>三个词的首字母缩写。

本节详细学习一下GVK相关的概念。我们在描述Kubernetes API时经常会用到这样一个四元组：<mark>Groups、Versions、Kinds和Resources</mark>。它们具体是什么含义呢？

- Groups和Versions：<u>一个Kubernetes API Group表示的是一些相关功能的集合</u>，比如apps这个Group里面就包含deployments、replicasets、daemonsets、statefulsets等资源，这些资源都是应用工作负载相关的，也就放在了同一个Group下。一个Group可以有一个或多个Versions，不难理解这里的用意，毕竟随着时间的推移，一个Group中的API难免有所变化。也许大家已经注意到，以前的Kubernetes版本创建Deployment时apiVersion用过apps/v1beta1和apps/v1beta2，现在已经是apps/v1了。

- Kinds和Resources：<mark>每个group-version（确定版本的一个组）中都包含一个或多个API类型，这些类型就是这里说的Kinds</mark>。每个Kind在不同的版本中一般会有所差异，但是每个版本的Kind要能够存储其他版本Kind的资源类型，无论是通过存储在字段里实现还是通过存储在注解中实现，具体的机制后面会详细讲解。这也就意味着使用老版本的API存储新版本类型数据不会引起数据丢失或污染。<mark>至于Resources，指的是一个Kind的具体使用，比如Pod类型对应的资源是pods。</mark>例如我们可以创建5个pods资源，其类型是Pod。描述资源的单词都是小写的，就像pods，而对应的类型一般就是这个资源的首字母大写单词的单数形式，比如pods对应Pod。类型和资源往往是一一对应的，尤其是在CRD的实现上。常见的特例就是为了支持HorizontalPodAutoscaler(HPA)和不同类型交互，Scale类型对应的资源有deployments/scale和replicasets/scale两种。

于是我们知道了可以<mark>通过一个GroupVersionKind(GVK)确定一个具体的类型</mark>，同样<mark>确定一个资源也就可以通过GroupVersionResource(GVR)来实现</mark>。
