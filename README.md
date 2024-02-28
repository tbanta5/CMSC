# CMSC Coffee-API
As a collaborative capstone project, the developers (group 2) of the Coffee-Api aim to create a general and easily extensible API which can be used by any Coffee Shop along with their own customized front end application. 

That being said, a comfortable development workflow was established following microservice principles using ![Docker](https://www.docker.com/products/docker-desktop/) and ![Kubernetes in Docker (KiND)](https://kind.sigs.k8s.io/).

## Design

![alt text](https://github.com/tbanta5/CMSC/blob/main/.readme-images/Coffee_NO_Java_Schema.png)

## Development
Development within this environment is meant to be repeatable, cohesive and transparent. To achieve this, we use Makefile to represent all patterns within development and as the source of truth for the project. This informs the user of what will be done and boosts confidence in the underlying processes that would otherwise be "magic".

### Setup
To setup this project in your own environment, ensure the Tools listed below are downloaded first according to Operating System. 
##### Tools for Development Environment - Mac Users 🍎.
Tools leveraged to create a cohesive and approchable development environment are:
- ![VSCode](https://code.visualstudio.com/download)
- ![Golang](https://go.dev/dl/)
- ![Docker](https://www.docker.com/products/docker-desktop/)
- ![Homebrew Package Manager](https://brew.sh/)

##### Tools for Development Environment - Windows Users🪟🪟.
- ![VSCode](https://code.visualstudio.com/download)
- ![Golang](https://go.dev/dl/)
- ![Chocolatey Package Manager](https://chocolatey.org/install)
- ![Beta - Homebrew Package Manager](https://brew.sh/)

#### Using Makefile
*** NOTE: For Windows users, you can use either Chocolatey or Homebrew (currently untested). Ensure you run all chocolatey commands from a privileged command prompt. Windows users must first install ![Make](https://community.chocolatey.org/packages/make) by running the below:
`'choco install make'`

Once the foundational tools are on your system, you can examine the Makefile commands and finish the enviroment setup by running the make command below based on your environment:
`make setup.[mac | windows]`

At this point the development environment should be configured and ready for use.

Let's test it out by running `make kind-up` which will create a new kubernetes cluster in docker to send our application to. Running `make kind-down` will delete the kubernetes cluster.

### Development Patterns
`make build` Builds all changes made to the application from /cmd and /internal directories (or subdirectories) into a new container image. 

`make kind-load` Loads the newly built container into Kubernetes environment.

`make kind-apply` Deploys all Kubernetes infrastructure defined in the /k8s directory. (ie. The application which was built and loaded as well as the database).

`make kind-delete` Destroys all deployed kubernetes objects, leaving an empty kind environment. ( Run `make kind-apply` to redeploy them.)

`make kind-restart` Replaces the old application with a newly built and loaded image but doesn't affect the rest of the kubernetes objects (like the database). 

`make kind-down` Destroys everything, including the kubernetes cluster.

`make kind-up` Creates a kubernetes cluster.

`make kind-status` Look at the status of deployed kubernetes resources.

`make kind-logs` Look at the logs of the coffee-api application from kubernetes.

## Educational Model
More than just establishing a development pattern, the hope is that this outline could model potential classroom lab configurations. 

### Why Container Technology
The pervasiveness, ease of use and accessibility of container technology coupled with the free allowances from ![Docker for Academic institutions](https://www.docker.com/community/open-source/application/) make for a perfect way to distribute software to students with low overhead, little static and expedited usage without the hassle of setup. 

`Hassle` has been the unfortunate norm for all Java based labs at UMGC. With experienced engineers spending upwards of 3 hours just to setup the environment for classwork. To put in current software terms, the Java Labs have a negative impact on developer "Velocity".

Containers are the standard way to distribute software in professional settings because they facilitate reproducable development environments. Containers themselves are based on Linux Operating system fundamentals which are seldom taugh in acadamia although widely used. 

Teaching container technology fundamentals along with Linux/Unix concepts (ie common Bash commands) is paramount to institutions that advertise "current technologies and job readiness" to students. The good news is that this material could be covered in 1 introductory class. Said class also serving to guide students in the setup and basic usage of a container environment they will continually use and build upon for all subsequent classes. 

### Microservices
The microservice pattern is heavily used in current job markets and builds on container technology fundamentals. It also encompases concepts of System Design - which is extremely important in real world software development. 

A microservice environment can easily grow and expand with students as they learn new concepts and materials. A student could learn to build web applications that interface with a docker contained database or perhaps said student in a later course is focusing on details of authentication/authorization in a full web architecture running in Kubernetes. 

Following configuration as code best practices in kubernetes manifests allows for a reproducable deployment environment among all students with any given set of applications. This allows educators more time for "teaching" rather than helpdesk work and cirriculum to be hyper focused on subject mastery (without concerns of complexities needed for setup).

### Core Language Selection - Golang
 A single language would be taught throughout the life of the university cirriculum with growing complexities in and around it. Concepts like infrastructure, concurrency, parallelism, web architectures, data structures, algorithms, and security could all help to expand and deepen student mastery of the language.

Language selection is important. Golang is ubiquitous in container technology and dominates the majority of ![cloud native projects](https://jonathonhenderson.co.uk/2023/07/16/cncf-projects-by-language). Go aims to be simple and performant and comes with a ![compatibility promise](https://go.dev/blog/compat) that all software will be backwards compatible through the life of Go v1. From a maintainability standpoint, this means the code written today will still run in 5-10 years and can be easily updated with automated CI/CD systems.

Many users find Golang to be easy to read and write (similar to python) with performance speeds comparable to C++ and C. It is a general purpose language built with systems programming in mind and natively supports cross compilation for any OS architecture. This is why Go is also found in embedded architectures like Kubernetes and Docker. Go additionally has a wealth of tooling and language support in current development environments like VSCode. Conversely, the complexities that come with the Java JVM, resource heavy development environments like Eclipse and Netbeans, and sheer size of the Java language make it less user friendly, more complicated to setup, more resource intensive with no performance benefits when compared to go. 

 Why not adopt a language that is easier to read and write, just as performant and ubiquitous in real world market to put students on the leading edge of technology? 
