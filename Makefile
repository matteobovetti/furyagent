

build:
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' .
	mv furyctl bin

copy-and-test-furyctl:
	sudo docker-compose  up -d 
	sudo docker cp bin/furyctl furyctl_etcd_1:/bin
	sudo docker cp furyagent.yml furyctl_etcd_1:/
	sudo docker cp test.sh furyctl_etcd_1:/
	sudo docker exec -ti furyctl_etcd_1 sh -c "mkdir -p /etc/etcd/pki && touch /etc/etcd/pki/ca.crt /etc/etcd/pki/ca.key"
	sudo docker exec -ti furyctl_etcd_1 sh -c "chmod u+x test.sh && ./test.sh"
	sudo docker exec -ti furyctl_etcd_1 usr/local/bin/etcd --name s1 --data-dir /etcd-data --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2379 --listen-peer-urls http://0.0.0.0:2380 --initial-advertise-peer-urls http://0.0.0.0:2380 --initial-cluster s1=http://0.0.0.0:2380 --initial-cluster-token tkn --initial-cluster-state new 1>/dev/null
	sudo docker exec -ti furyctl_etcd_1 sh -c "echo 'read after restart' && etcdctl get foo"
	#sudo docker-compose logs -f

vendor:
	go mod vendor

install_deps:
	go get github.com/mitchellh/gox
	
up:
	docker-compose up -d && docker-compose logs -f

down: 
	docker-compose down -v
