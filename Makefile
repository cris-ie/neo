.PHONY: build deploy build

export KUBECONFIG=${HOME}/.kube/config_k3s

build:
	docker build -t neo .

run: build
	docker run -it --rm neo 

deploy:
	helm dep update chart
	helm upgrade -n neo -i neo chart --atomic --wait --timeout 200s --create-namespace --debug

clean:
	vagrant destroy -f

vagrant:
	vagrant up
	scp -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i ./.vagrant/machines/k3s/virtualbox/private_key -r -P2222 vagrant@127.0.0.1:/home/vagrant/.kube/config_external ~/.kube/config_k3s
	@echo "export KUBECONFIG=~/.kube/config_k3s"

full: vagrant deploy
	curl http://neo.192.168.33.10.nip.io/status
	curl http://neo.192.168.33.10.nip.io/liveness
	curl http://neo.192.168.33.10.nip.io/neo/week
	curl http://neo.192.168.33.10.nip.io/neo/next
