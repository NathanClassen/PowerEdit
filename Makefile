build:
	go build -o bin/poweredit


install:
	cp bin/poweredit ~/go/bin
	
clean:
	rm bin/poweredit