# Klouds.io / Klouds.org

#####We like to leverage automation:

[![wercker status](https://app.wercker.com/status/fcf719ade20c4ab01184d966f4650ee2/s/master "wercker status")](https://app.wercker.com/project/bykey/fcf719ade20c4ab01184d966f4650ee2)  [![Join the chat at https://gitter.im/Superordinate/klouds](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/Superordinate/klouds?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)  [![Code Climate](https://codeclimate.com/github/Superordinate/klouds/badges/gpa.svg)](https://codeclimate.com/github/Superordinate/klouds)  


<img src="http://www.ozzadar.com/klouds.png" align="center"/>



A web framework for dynamically launching and accessing web servers.

Uses rest APIs to manage containers. Uses wercker to automate deployment.  (maybe)


##THE STACK

```

mesos-slave/consul --\   /--------------------------------------> haproxy-consul  <---\
			   		  \	/															   \
mesos-slave/consul ---->mesos-master ->- marathon endpoint --><-- klouds-frontend <-----User
			   		  /	\					    v
mesos-slave/consul __/	 \______<_______ marathon-consul

```


##HOW TO USE:

###Environment Variables

MARATHON_ENDPOINT= 192.168.3.4:8080  	<-- Point this to your own marathon backend
KLOUDS_DOMAIN=klouds.org 				<-- Creates application access urls to point to your domain
MYSQL_HOST= 127.0.0.1:3306	 			<-- Points to your database
MYSQL_USER= root						<-- User for your mysql database
MYSQL_PASSWORD= iamapassword			<-- password for you mysql user


###To build (linux):


```
go get github.com/superordinate/klouds
cd $GOPATH/src/github.com/superordinate/klouds
go build .

```
### To Run

``` 

	MYSQL_HOST=127.0.0.1:3306 MYSQL_USER=root MYSQL_PASSWORD=iamapassword MARATHON_ENDPOINT=192.168.3.4:8080 KLOUDS_DOMAIN=mydomain.com ./klouds

```

OR

```
	export KLOUDS_DOMAIN=mydomain.com 
	export MYSQL_HOST= 127.0.0.1:3306
	export MYSQL_USER= root	
	export MYSQL_PASSWORD= iamapassword
	export MARATHON_ENDPOINT=192.168.3.4:8080

	./klouds

```


<a href="https://github.com/EmileVauge/traefik"><img src="https://camo.githubusercontent.com/0d83f4ec95b28ecc0353078ca4364bf461b99c2d/687474703a2f2f7472616566696b2e6769746875622e696f2f7472616566696b2e6c6f676f2e737667" align="center" height="96" width="200" ></a><br>

