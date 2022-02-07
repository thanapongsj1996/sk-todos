## Developing Robust API Services with GO

**Course link :** https://www.skooldio.com/courses/developing-robust-api-services-with-go <br>
**Build command :** ``go build -ldflags "-X main.buildcommit=`git rev-parse --short HEAD` -X main.buildtime=`date "+%Y-%m-%dT%H:%M:%S%Z:00"`" -o app`` <br>
**Certificate :** `https://www.skooldio.com/certificate/a6395f9c-fb71-43f9-b68a-723d62a7d979`<br>
**Load test :** `echo "GET http://localhost:8080/load-test" | vegeta attack -rate=30 -duration=3s | vegeta report` <br>
**Load test with result :** `echo "GET http://localhost:8080/ping" | vegeta attack -name=20qps -rate=50 -duration=3s > results.bin` <br>
**Plot results from load test :** `vegeta plot results.bin > plot.html`
