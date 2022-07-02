# 需求
## custom controller需要监听service资源变化，当service资源变化时:
## 如果新增Service时
## 1.包含指定annotation，创建ingress对象
##   不包含指定annotation,忽略
## 2.如果删除service，删除ingress
## 3.如果更新service
##   指定包含annotation，检查ingress资源对象是否存在，如果不存在，则建立，存在则忽略

# Ingress部署
## 1.参考文档https://kubernetes.github.io/ingress-nginx/user-guide/basic-usage/
