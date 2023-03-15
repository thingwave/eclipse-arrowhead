
all:
	go build -o ah-datamanager

all-arm64:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ah-datamanager

run:
	./ah-datamanager

clean:
	rm -rf ah-datamanager

gtest:
	curl --insecure -X GET https://127.0.0.1:8461/dm/historian/Mulle-342/temp

ptest:
	curl --insecure -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle-342","bu":"Cel","ver":5,"n":"temperature","v":21.2}]' https://127.0.0.1:8461/datamanager/proxy/Mulle-342/temp

putest:
	curl -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle-342","bu":"Cel","ver":5,"n":"temperature","v":20.1}]' http://127.0.0.1:8461/datamanager/proxy/mulle-342/temp

hutest:
	curl -H 'Content-Type: application/json' -X PUT -d '[{"bn":"mulle-342","bu":"Cel","ver":5,"n":"temperature","v":21.2}]' http://127.0.0.1:8461/datamanager/historian/mulle-342/temp




