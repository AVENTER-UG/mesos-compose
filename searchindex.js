Object.assign(window.search, {"doc_urls":["index.html#introduction---mesos-compose-the-docker-compose-framework-for-apache-mesos","index.html#requirements","index.html#example","errorcode.html#error-codes","labels.html#m3s-labels","labels.html#bizaventermesos_composecontainer_type","labels.html#bizaventermesos_composecontraint_hostname"],"index":{"documentStore":{"docInfo":{"0":{"body":0,"breadcrumbs":9,"title":8},"1":{"body":14,"breadcrumbs":2,"title":1},"2":{"body":86,"breadcrumbs":2,"title":1},"3":{"body":12,"breadcrumbs":3,"title":2},"4":{"body":0,"breadcrumbs":5,"title":2},"5":{"body":12,"breadcrumbs":4,"title":1},"6":{"body":6,"breadcrumbs":4,"title":1}},"docs":{"0":{"body":"","breadcrumbs":"Introduction » Introduction - mesos-compose, the docker-compose framework for Apache Mesos","id":"0","title":"Introduction - mesos-compose, the docker-compose framework for Apache Mesos"},"1":{"body":"Apache Mesos min 1.6.0 Mesos with SSL and Authentication is optional Redis Database Docker Compose Spec 3.9","breadcrumbs":"Introduction » Requirements","id":"1","title":"Requirements"},"2":{"body":"The compose file: version: '3.9'\nservices: app: image: alpine:latest command: [\"sleep\",\"10\"] restart: always volumes: - /tmp:/tmp environment: - MYSQL_HOST=test labels: biz.aventer.mesos_compose.container_type: \"DOCKER\" network_mode: \"BRIDGE\" network: - default deploy: resources: limits: cpus: \"0.001\" memory: \"50\" networks: default: external: true name: weave Push these compose file to the framework. Every compose file needs to have an own project name. curl -X PUT http://localhost:10000/v0/compose/<PROJECTNAME> --data-binary @docs/example/docker-compose.yml image_2021-11-08-11-33-09 image_2021-11-08-11-33-47 To scale the service, just execute the same call again. To update a already existing docker-compose project, call: curl -X PUT http://localhost:10000/v0/compose/<PROJECTNAME>/update --data-binary @docs/example/docker-compose.yml","breadcrumbs":"Introduction » Example","id":"2","title":"Example"},"3":{"body":"| 0 | No Error | | 1 | Could not save configuration into Redis | | 2 | Unspecified error | | 3 | Reconcile Error |","breadcrumbs":"Errorcodes » Error Codes","id":"3","title":"Error Codes"},"4":{"body":"","breadcrumbs":"Custom M3s Labels » M3s Labels","id":"4","title":"M3s Labels"},"5":{"body":"Values: docker, mesos (todo) This label will control, which container engine mesos will choose to execute the container.","breadcrumbs":"Custom M3s Labels » biz.aventer.mesos_compose.container_type","id":"5","title":"biz.aventer.mesos_compose.container_type"},"6":{"body":"Values: This label will control, on which node the container will be executed.","breadcrumbs":"Custom M3s Labels » biz.aventer.mesos_compose.contraint_hostname","id":"6","title":"biz.aventer.mesos_compose.contraint_hostname"}},"length":7,"save":true},"fields":["title","body","breadcrumbs"],"index":{"body":{"root":{"0":{".":{"0":{"0":{"1":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"8":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"9":{"df":1,"docs":{"2":{"tf":1.0}}},"df":1,"docs":{"3":{"tf":1.0}}},"1":{".":{"6":{".":{"0":{"df":1,"docs":{"1":{"tf":1.0}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"1":{"df":1,"docs":{"2":{"tf":2.0}}},"df":1,"docs":{"3":{"tf":1.0}}},"2":{"df":1,"docs":{"3":{"tf":1.0}}},"3":{".":{"9":{"df":2,"docs":{"1":{"tf":1.0},"2":{"tf":1.0}}},"df":0,"docs":{}},"3":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":1,"docs":{"3":{"tf":1.0}}},"4":{"7":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"5":{"0":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"a":{"df":0,"docs":{},"g":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}},"df":0,"docs":{}},"l":{"df":0,"docs":{},"p":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"e":{":":{"df":0,"docs":{},"l":{"a":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"df":0,"docs":{}}},"df":0,"docs":{}}}}},"r":{"df":0,"docs":{},"e":{"a":{"d":{"df":0,"docs":{},"i":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"w":{"a":{"df":0,"docs":{},"y":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}}},"p":{"a":{"c":{"df":0,"docs":{},"h":{"df":2,"docs":{"0":{"tf":1.0},"1":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{},"p":{"df":1,"docs":{"2":{"tf":1.0}}}},"u":{"df":0,"docs":{},"t":{"df":0,"docs":{},"h":{"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"df":1,"docs":{"1":{"tf":1.0}}}}}}}}},"b":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"r":{"df":0,"docs":{},"i":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"df":0,"docs":{}},"z":{".":{"a":{"df":0,"docs":{},"v":{"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{".":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"_":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"e":{".":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"_":{"df":0,"docs":{},"t":{"df":0,"docs":{},"y":{"df":0,"docs":{},"p":{"df":2,"docs":{"2":{"tf":1.0},"5":{"tf":1.0}}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{},"r":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"_":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"m":{"df":1,"docs":{"6":{"tf":1.0}}}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"r":{"df":0,"docs":{},"i":{"d":{"df":0,"docs":{},"g":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}}}},"c":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":1,"docs":{"5":{"tf":1.0}}}}}},"o":{"d":{"df":0,"docs":{},"e":{"df":1,"docs":{"3":{"tf":1.0}}}},"df":0,"docs":{},"m":{"df":0,"docs":{},"m":{"a":{"df":0,"docs":{},"n":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}},"df":0,"docs":{}},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":3,"docs":{"0":{"tf":1.4142135623730951},"1":{"tf":1.0},"2":{"tf":2.0}},"e":{".":{"df":0,"docs":{},"y":{"df":0,"docs":{},"m":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}},"df":0,"docs":{}}}}}},"n":{"df":0,"docs":{},"f":{"df":0,"docs":{},"i":{"df":0,"docs":{},"g":{"df":0,"docs":{},"u":{"df":0,"docs":{},"r":{"df":1,"docs":{"3":{"tf":1.0}}}}}}},"t":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":2,"docs":{"5":{"tf":1.4142135623730951},"6":{"tf":1.0}}}}},"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"l":{"df":2,"docs":{"5":{"tf":1.0},"6":{"tf":1.0}}}}}}}},"p":{"df":0,"docs":{},"u":{"df":1,"docs":{"2":{"tf":1.0}}}},"u":{"df":0,"docs":{},"r":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}},"d":{"a":{"df":0,"docs":{},"t":{"a":{"b":{"a":{"df":0,"docs":{},"s":{"df":1,"docs":{"1":{"tf":1.0}}}},"df":0,"docs":{}},"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":0,"docs":{}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"f":{"a":{"df":0,"docs":{},"u":{"df":0,"docs":{},"l":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}},"df":0,"docs":{}},"p":{"df":0,"docs":{},"l":{"df":0,"docs":{},"o":{"df":0,"docs":{},"y":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"o":{"c":{"df":0,"docs":{},"k":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":4,"docs":{"0":{"tf":1.0},"1":{"tf":1.0},"2":{"tf":1.4142135623730951},"5":{"tf":1.0}}}}},"s":{"/":{"df":0,"docs":{},"e":{"df":0,"docs":{},"x":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"l":{"df":0,"docs":{},"e":{"/":{"d":{"df":0,"docs":{},"o":{"c":{"df":0,"docs":{},"k":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}},"df":0,"docs":{}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}},"df":0,"docs":{}}}},"df":0,"docs":{}}},"df":0,"docs":{}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"g":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":1,"docs":{"5":{"tf":1.0}}}}},"v":{"df":0,"docs":{},"i":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}}}}},"r":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":1,"docs":{"3":{"tf":2.0}}}}}},"x":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"u":{"df":0,"docs":{},"t":{"df":3,"docs":{"2":{"tf":1.0},"5":{"tf":1.0},"6":{"tf":1.0}}}}},"df":0,"docs":{}},"i":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}}}}},"f":{"df":0,"docs":{},"i":{"df":0,"docs":{},"l":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.7320508075688772}}}}},"r":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":0,"docs":{},"w":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"k":{"df":2,"docs":{"0":{"tf":1.0},"2":{"tf":1.0}}}}}}}}},"df":0,"docs":{}}},"h":{"df":0,"docs":{},"t":{"df":0,"docs":{},"t":{"df":0,"docs":{},"p":{":":{"/":{"/":{"df":0,"docs":{},"l":{"df":0,"docs":{},"o":{"c":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{":":{"1":{"0":{"0":{"0":{"0":{"/":{"df":0,"docs":{},"v":{"0":{"/":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"e":{"/":{"<":{"df":0,"docs":{},"p":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"j":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"t":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"m":{"df":1,"docs":{"2":{"tf":1.0}},"e":{">":{"/":{"df":0,"docs":{},"u":{"df":0,"docs":{},"p":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}}},"df":0,"docs":{}},"df":0,"docs":{}}}},"df":0,"docs":{}}}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}}},"i":{"df":0,"docs":{},"m":{"a":{"df":0,"docs":{},"g":{"df":1,"docs":{"2":{"tf":1.0}},"e":{"_":{"2":{"0":{"2":{"1":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}},"df":0,"docs":{}},"n":{"df":0,"docs":{},"t":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"d":{"df":0,"docs":{},"u":{"c":{"df":0,"docs":{},"t":{"df":1,"docs":{"0":{"tf":1.0}}}},"df":0,"docs":{}}},"df":0,"docs":{}}}}}},"l":{"a":{"b":{"df":0,"docs":{},"e":{"df":0,"docs":{},"l":{"df":4,"docs":{"2":{"tf":1.0},"4":{"tf":1.0},"5":{"tf":1.0},"6":{"tf":1.0}}}}},"df":0,"docs":{}},"df":0,"docs":{},"i":{"df":0,"docs":{},"m":{"df":0,"docs":{},"i":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"m":{"3":{"df":1,"docs":{"4":{"tf":1.0}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"m":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"i":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"s":{"df":0,"docs":{},"o":{"df":3,"docs":{"0":{"tf":1.4142135623730951},"1":{"tf":1.4142135623730951},"5":{"tf":1.4142135623730951}}}}},"i":{"df":0,"docs":{},"n":{"df":1,"docs":{"1":{"tf":1.0}}}},"y":{"df":0,"docs":{},"s":{"df":0,"docs":{},"q":{"df":0,"docs":{},"l":{"_":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"=":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{}}}}}},"n":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"e":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"t":{"df":0,"docs":{},"w":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"k":{"_":{"df":0,"docs":{},"m":{"df":0,"docs":{},"o":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}}},"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}}}},"o":{"d":{"df":0,"docs":{},"e":{"df":1,"docs":{"6":{"tf":1.0}}}},"df":0,"docs":{}}},"o":{"df":0,"docs":{},"p":{"df":0,"docs":{},"t":{"df":0,"docs":{},"i":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":1,"docs":{"1":{"tf":1.0}}}}}}}},"p":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"j":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}},"df":0,"docs":{}}}}},"u":{"df":0,"docs":{},"s":{"df":0,"docs":{},"h":{"df":1,"docs":{"2":{"tf":1.0}}}},"t":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"r":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"c":{"df":0,"docs":{},"i":{"df":0,"docs":{},"l":{"df":1,"docs":{"3":{"tf":1.0}}}}},"df":0,"docs":{}}}},"d":{"df":0,"docs":{},"i":{"df":2,"docs":{"1":{"tf":1.0},"3":{"tf":1.0}}}},"df":0,"docs":{},"q":{"df":0,"docs":{},"u":{"df":0,"docs":{},"i":{"df":0,"docs":{},"r":{"df":1,"docs":{"1":{"tf":1.0}}}}}},"s":{"df":0,"docs":{},"o":{"df":0,"docs":{},"u":{"df":0,"docs":{},"r":{"c":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}}},"t":{"a":{"df":0,"docs":{},"r":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}},"df":0,"docs":{}}}}},"s":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.0}}}},"v":{"df":0,"docs":{},"e":{"df":1,"docs":{"3":{"tf":1.0}}}}},"c":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.0}}}}},"df":0,"docs":{}},"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"v":{"df":0,"docs":{},"i":{"c":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":0,"docs":{}}}}},"l":{"df":0,"docs":{},"e":{"df":0,"docs":{},"e":{"df":0,"docs":{},"p":{"\"":{",":{"\"":{"1":{"0":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}}},"p":{"df":0,"docs":{},"e":{"c":{"df":1,"docs":{"1":{"tf":1.0}}},"df":0,"docs":{}}},"s":{"df":0,"docs":{},"l":{"df":1,"docs":{"1":{"tf":1.0}}}}},"t":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{":":{"/":{"df":0,"docs":{},"t":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"o":{"d":{"df":0,"docs":{},"o":{"df":1,"docs":{"5":{"tf":1.0}}}},"df":0,"docs":{}},"r":{"df":0,"docs":{},"u":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"u":{"df":0,"docs":{},"n":{"df":0,"docs":{},"s":{"df":0,"docs":{},"p":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"i":{"df":0,"docs":{},"f":{"df":0,"docs":{},"i":{"df":1,"docs":{"3":{"tf":1.0}}}}}},"df":0,"docs":{}}}}},"p":{"d":{"a":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"v":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"u":{"df":2,"docs":{"5":{"tf":1.0},"6":{"tf":1.0}}}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"s":{"df":0,"docs":{},"i":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}}}}},"o":{"df":0,"docs":{},"l":{"df":0,"docs":{},"u":{"df":0,"docs":{},"m":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"w":{"df":0,"docs":{},"e":{"a":{"df":0,"docs":{},"v":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}}},"x":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"breadcrumbs":{"root":{"0":{".":{"0":{"0":{"1":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"8":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"9":{"df":1,"docs":{"2":{"tf":1.0}}},"df":1,"docs":{"3":{"tf":1.0}}},"1":{".":{"6":{".":{"0":{"df":1,"docs":{"1":{"tf":1.0}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"1":{"df":1,"docs":{"2":{"tf":2.0}}},"df":1,"docs":{"3":{"tf":1.0}}},"2":{"df":1,"docs":{"3":{"tf":1.0}}},"3":{".":{"9":{"df":2,"docs":{"1":{"tf":1.0},"2":{"tf":1.0}}},"df":0,"docs":{}},"3":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":1,"docs":{"3":{"tf":1.0}}},"4":{"7":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"5":{"0":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"a":{"df":0,"docs":{},"g":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}},"df":0,"docs":{}},"l":{"df":0,"docs":{},"p":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"e":{":":{"df":0,"docs":{},"l":{"a":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"df":0,"docs":{}}},"df":0,"docs":{}}}}},"r":{"df":0,"docs":{},"e":{"a":{"d":{"df":0,"docs":{},"i":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"w":{"a":{"df":0,"docs":{},"y":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}}},"p":{"a":{"c":{"df":0,"docs":{},"h":{"df":2,"docs":{"0":{"tf":1.4142135623730951},"1":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{},"p":{"df":1,"docs":{"2":{"tf":1.0}}}},"u":{"df":0,"docs":{},"t":{"df":0,"docs":{},"h":{"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"df":1,"docs":{"1":{"tf":1.0}}}}}}}}},"b":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"r":{"df":0,"docs":{},"i":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"df":0,"docs":{}},"z":{".":{"a":{"df":0,"docs":{},"v":{"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{".":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"_":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"e":{".":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"_":{"df":0,"docs":{},"t":{"df":0,"docs":{},"y":{"df":0,"docs":{},"p":{"df":2,"docs":{"2":{"tf":1.0},"5":{"tf":1.4142135623730951}}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{},"r":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"_":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"m":{"df":1,"docs":{"6":{"tf":1.4142135623730951}}}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"r":{"df":0,"docs":{},"i":{"d":{"df":0,"docs":{},"g":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}}}},"c":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":1,"docs":{"5":{"tf":1.0}}}}}},"o":{"d":{"df":0,"docs":{},"e":{"df":1,"docs":{"3":{"tf":1.4142135623730951}}}},"df":0,"docs":{},"m":{"df":0,"docs":{},"m":{"a":{"df":0,"docs":{},"n":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}},"df":0,"docs":{}},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":3,"docs":{"0":{"tf":2.0},"1":{"tf":1.0},"2":{"tf":2.0}},"e":{".":{"df":0,"docs":{},"y":{"df":0,"docs":{},"m":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}},"df":0,"docs":{}}}}}},"n":{"df":0,"docs":{},"f":{"df":0,"docs":{},"i":{"df":0,"docs":{},"g":{"df":0,"docs":{},"u":{"df":0,"docs":{},"r":{"df":1,"docs":{"3":{"tf":1.0}}}}}}},"t":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":2,"docs":{"5":{"tf":1.4142135623730951},"6":{"tf":1.0}}}}},"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"l":{"df":2,"docs":{"5":{"tf":1.0},"6":{"tf":1.0}}}}}}}},"p":{"df":0,"docs":{},"u":{"df":1,"docs":{"2":{"tf":1.0}}}},"u":{"df":0,"docs":{},"r":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}},"s":{"df":0,"docs":{},"t":{"df":0,"docs":{},"o":{"df":0,"docs":{},"m":{"df":3,"docs":{"4":{"tf":1.0},"5":{"tf":1.0},"6":{"tf":1.0}}}}}}}},"d":{"a":{"df":0,"docs":{},"t":{"a":{"b":{"a":{"df":0,"docs":{},"s":{"df":1,"docs":{"1":{"tf":1.0}}}},"df":0,"docs":{}},"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":0,"docs":{}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"f":{"a":{"df":0,"docs":{},"u":{"df":0,"docs":{},"l":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}},"df":0,"docs":{}},"p":{"df":0,"docs":{},"l":{"df":0,"docs":{},"o":{"df":0,"docs":{},"y":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"o":{"c":{"df":0,"docs":{},"k":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":4,"docs":{"0":{"tf":1.4142135623730951},"1":{"tf":1.0},"2":{"tf":1.4142135623730951},"5":{"tf":1.0}}}}},"s":{"/":{"df":0,"docs":{},"e":{"df":0,"docs":{},"x":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"l":{"df":0,"docs":{},"e":{"/":{"d":{"df":0,"docs":{},"o":{"c":{"df":0,"docs":{},"k":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}},"df":0,"docs":{}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}},"df":0,"docs":{}}}},"df":0,"docs":{}}},"df":0,"docs":{}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"g":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":1,"docs":{"5":{"tf":1.0}}}}},"v":{"df":0,"docs":{},"i":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}}}}},"r":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"c":{"df":0,"docs":{},"o":{"d":{"df":1,"docs":{"3":{"tf":1.0}}},"df":0,"docs":{}}},"df":1,"docs":{"3":{"tf":2.23606797749979}}}}}},"x":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}},"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"u":{"df":0,"docs":{},"t":{"df":3,"docs":{"2":{"tf":1.0},"5":{"tf":1.0},"6":{"tf":1.0}}}}},"df":0,"docs":{}},"i":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}}}}},"f":{"df":0,"docs":{},"i":{"df":0,"docs":{},"l":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.7320508075688772}}}}},"r":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":0,"docs":{},"w":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"k":{"df":2,"docs":{"0":{"tf":1.4142135623730951},"2":{"tf":1.0}}}}}}}}},"df":0,"docs":{}}},"h":{"df":0,"docs":{},"t":{"df":0,"docs":{},"t":{"df":0,"docs":{},"p":{":":{"/":{"/":{"df":0,"docs":{},"l":{"df":0,"docs":{},"o":{"c":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{":":{"1":{"0":{"0":{"0":{"0":{"/":{"df":0,"docs":{},"v":{"0":{"/":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"e":{"/":{"<":{"df":0,"docs":{},"p":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"j":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"t":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"m":{"df":1,"docs":{"2":{"tf":1.0}},"e":{">":{"/":{"df":0,"docs":{},"u":{"df":0,"docs":{},"p":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}}},"df":0,"docs":{}},"df":0,"docs":{}}}},"df":0,"docs":{}}}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}}},"i":{"df":0,"docs":{},"m":{"a":{"df":0,"docs":{},"g":{"df":1,"docs":{"2":{"tf":1.0}},"e":{"_":{"2":{"0":{"2":{"1":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}},"df":0,"docs":{}},"n":{"df":0,"docs":{},"t":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"d":{"df":0,"docs":{},"u":{"c":{"df":0,"docs":{},"t":{"df":3,"docs":{"0":{"tf":1.7320508075688772},"1":{"tf":1.0},"2":{"tf":1.0}}}},"df":0,"docs":{}}},"df":0,"docs":{}}}}}},"l":{"a":{"b":{"df":0,"docs":{},"e":{"df":0,"docs":{},"l":{"df":4,"docs":{"2":{"tf":1.0},"4":{"tf":1.7320508075688772},"5":{"tf":1.4142135623730951},"6":{"tf":1.4142135623730951}}}}},"df":0,"docs":{}},"df":0,"docs":{},"i":{"df":0,"docs":{},"m":{"df":0,"docs":{},"i":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"m":{"3":{"df":3,"docs":{"4":{"tf":1.7320508075688772},"5":{"tf":1.0},"6":{"tf":1.0}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"m":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"i":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"s":{"df":0,"docs":{},"o":{"df":3,"docs":{"0":{"tf":2.0},"1":{"tf":1.4142135623730951},"5":{"tf":1.4142135623730951}}}}},"i":{"df":0,"docs":{},"n":{"df":1,"docs":{"1":{"tf":1.0}}}},"y":{"df":0,"docs":{},"s":{"df":0,"docs":{},"q":{"df":0,"docs":{},"l":{"_":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"=":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{}}}}}},"n":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"e":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"t":{"df":0,"docs":{},"w":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"k":{"_":{"df":0,"docs":{},"m":{"df":0,"docs":{},"o":{"d":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}}},"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}}}}},"o":{"d":{"df":0,"docs":{},"e":{"df":1,"docs":{"6":{"tf":1.0}}}},"df":0,"docs":{}}},"o":{"df":0,"docs":{},"p":{"df":0,"docs":{},"t":{"df":0,"docs":{},"i":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":1,"docs":{"1":{"tf":1.0}}}}}}}},"p":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"j":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}},"df":0,"docs":{}}}}},"u":{"df":0,"docs":{},"s":{"df":0,"docs":{},"h":{"df":1,"docs":{"2":{"tf":1.0}}}},"t":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"r":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"c":{"df":0,"docs":{},"i":{"df":0,"docs":{},"l":{"df":1,"docs":{"3":{"tf":1.0}}}}},"df":0,"docs":{}}}},"d":{"df":0,"docs":{},"i":{"df":2,"docs":{"1":{"tf":1.0},"3":{"tf":1.0}}}},"df":0,"docs":{},"q":{"df":0,"docs":{},"u":{"df":0,"docs":{},"i":{"df":0,"docs":{},"r":{"df":1,"docs":{"1":{"tf":1.4142135623730951}}}}}},"s":{"df":0,"docs":{},"o":{"df":0,"docs":{},"u":{"df":0,"docs":{},"r":{"c":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}}}},"t":{"a":{"df":0,"docs":{},"r":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}}},"df":0,"docs":{}}}}},"s":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.0}}}},"v":{"df":0,"docs":{},"e":{"df":1,"docs":{"3":{"tf":1.0}}}}},"c":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.0}}}}},"df":0,"docs":{}},"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"v":{"df":0,"docs":{},"i":{"c":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}},"df":0,"docs":{}}}}},"l":{"df":0,"docs":{},"e":{"df":0,"docs":{},"e":{"df":0,"docs":{},"p":{"\"":{",":{"\"":{"1":{"0":{"df":1,"docs":{"2":{"tf":1.0}}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}},"df":0,"docs":{}}}}},"p":{"df":0,"docs":{},"e":{"c":{"df":1,"docs":{"1":{"tf":1.0}}},"df":0,"docs":{}}},"s":{"df":0,"docs":{},"l":{"df":1,"docs":{"1":{"tf":1.0}}}}},"t":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{":":{"/":{"df":0,"docs":{},"t":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"o":{"d":{"df":0,"docs":{},"o":{"df":1,"docs":{"5":{"tf":1.0}}}},"df":0,"docs":{}},"r":{"df":0,"docs":{},"u":{"df":0,"docs":{},"e":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"u":{"df":0,"docs":{},"n":{"df":0,"docs":{},"s":{"df":0,"docs":{},"p":{"df":0,"docs":{},"e":{"c":{"df":0,"docs":{},"i":{"df":0,"docs":{},"f":{"df":0,"docs":{},"i":{"df":1,"docs":{"3":{"tf":1.0}}}}}},"df":0,"docs":{}}}}},"p":{"d":{"a":{"df":0,"docs":{},"t":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"v":{"a":{"df":0,"docs":{},"l":{"df":0,"docs":{},"u":{"df":2,"docs":{"5":{"tf":1.0},"6":{"tf":1.0}}}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"s":{"df":0,"docs":{},"i":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":1,"docs":{"2":{"tf":1.0}}}}}}}},"o":{"df":0,"docs":{},"l":{"df":0,"docs":{},"u":{"df":0,"docs":{},"m":{"df":1,"docs":{"2":{"tf":1.0}}}}}}},"w":{"df":0,"docs":{},"e":{"a":{"df":0,"docs":{},"v":{"df":1,"docs":{"2":{"tf":1.0}}}},"df":0,"docs":{}}},"x":{"df":1,"docs":{"2":{"tf":1.4142135623730951}}}}},"title":{"root":{"a":{"df":0,"docs":{},"p":{"a":{"c":{"df":0,"docs":{},"h":{"df":1,"docs":{"0":{"tf":1.0}}}},"df":0,"docs":{}},"df":0,"docs":{}}},"b":{"df":0,"docs":{},"i":{"df":0,"docs":{},"z":{".":{"a":{"df":0,"docs":{},"v":{"df":0,"docs":{},"e":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{".":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"_":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"e":{".":{"c":{"df":0,"docs":{},"o":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"_":{"df":0,"docs":{},"t":{"df":0,"docs":{},"y":{"df":0,"docs":{},"p":{"df":1,"docs":{"5":{"tf":1.0}}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{},"r":{"a":{"df":0,"docs":{},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"_":{"df":0,"docs":{},"h":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":0,"docs":{},"t":{"df":0,"docs":{},"n":{"a":{"df":0,"docs":{},"m":{"df":1,"docs":{"6":{"tf":1.0}}}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}}}}},"df":0,"docs":{}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}}}}},"df":0,"docs":{}}}}}}}},"df":0,"docs":{}},"df":0,"docs":{}}}},"c":{"df":0,"docs":{},"o":{"d":{"df":0,"docs":{},"e":{"df":1,"docs":{"3":{"tf":1.0}}}},"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"o":{"df":0,"docs":{},"s":{"df":1,"docs":{"0":{"tf":1.4142135623730951}}}}}}}},"d":{"df":0,"docs":{},"o":{"c":{"df":0,"docs":{},"k":{"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":1,"docs":{"0":{"tf":1.0}}}}}},"df":0,"docs":{}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"r":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":1,"docs":{"3":{"tf":1.0}}}}}},"x":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"p":{"df":0,"docs":{},"l":{"df":1,"docs":{"2":{"tf":1.0}}}}}},"df":0,"docs":{}}},"f":{"df":0,"docs":{},"r":{"a":{"df":0,"docs":{},"m":{"df":0,"docs":{},"e":{"df":0,"docs":{},"w":{"df":0,"docs":{},"o":{"df":0,"docs":{},"r":{"df":0,"docs":{},"k":{"df":1,"docs":{"0":{"tf":1.0}}}}}}}}},"df":0,"docs":{}}},"i":{"df":0,"docs":{},"n":{"df":0,"docs":{},"t":{"df":0,"docs":{},"r":{"df":0,"docs":{},"o":{"d":{"df":0,"docs":{},"u":{"c":{"df":0,"docs":{},"t":{"df":1,"docs":{"0":{"tf":1.0}}}},"df":0,"docs":{}}},"df":0,"docs":{}}}}}},"l":{"a":{"b":{"df":0,"docs":{},"e":{"df":0,"docs":{},"l":{"df":1,"docs":{"4":{"tf":1.0}}}}},"df":0,"docs":{}},"df":0,"docs":{}},"m":{"3":{"df":1,"docs":{"4":{"tf":1.0}}},"df":0,"docs":{},"e":{"df":0,"docs":{},"s":{"df":0,"docs":{},"o":{"df":1,"docs":{"0":{"tf":1.4142135623730951}}}}}},"r":{"df":0,"docs":{},"e":{"df":0,"docs":{},"q":{"df":0,"docs":{},"u":{"df":0,"docs":{},"i":{"df":0,"docs":{},"r":{"df":1,"docs":{"1":{"tf":1.0}}}}}}}}}}},"lang":"English","pipeline":["trimmer","stopWordFilter","stemmer"],"ref":"id","version":"0.9.5"},"results_options":{"limit_results":30,"teaser_word_count":30},"search_options":{"bool":"OR","expand":true,"fields":{"body":{"boost":1},"breadcrumbs":{"boost":1},"title":{"boost":2}}}});