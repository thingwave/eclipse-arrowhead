#SRC = datamanager.go dmdatabase.go routes.go router.go handlers.go datamodels.go accesscontrol.go

all:
	go build -o ah-datamanager
#$(SRC)

run:
	./ah-datamanager

clean:
	rm -rf ah-datamanager

dependencies:
	go get github.com/BurntSushi/toml
	go get -u github.com/gorilla/mux
	go get github.com/go-sql-driver/mysql

gtest:
	curl --insecure -X GET https://127.0.0.1:8461/dm/historian/Mulle-342/temp

ptest:
	curl --insecure -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle-342","bu":"Cel","ver":5,"n":"temperature","v":21.2}]' https://127.0.0.1:8461/datamanager/proxy/Mulle-342/temp

putest:
	#curl -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle342","bu":"Cel","ver":4,"n":"temp", "v":12.1}]' http://127.0.0.1:8461/datamanager/proxy/mulle342/temp
	curl -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle-342","bu":"Cel","ver":5,"n":"temperature","v":20.1}]' http://127.0.0.1:8461/datamanager/proxy/mulle-342/temp

hutest:
	#curl -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle342","bu":"Cel","ver":4,"n":"temp", "v":12.1}]' http://127.0.0.1:8461/datamanager/proxy/mulle342/temp
	curl -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle-342","bu":"Cel","ver":5,"n":"temperature","v":21.2}]' http://127.0.0.1:8461/datamanager/historian/mulle-342/temp

ortest1:
	curl -X GET "http://127.0.0.1:8441/orchestrator/mgmt/store/all_top_priority?direction=DESC&page=27&sort_field=updatedAt&item_per_page=20"

ortest3:
	curl -H 'Content-Type: application/json' -X POST -d '{"priorityMap": {"id1":69, "id2": 6}}' http://127.0.0.1:8441/orchestrator/mgmt/store/modify_priorities
	

