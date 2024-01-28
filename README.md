# CMSC Coffee-API
As a collaborative capstone project, the developers (group 2) of the Coffee-Api aim to create a general and easily extensible API which can be used by any Coffee Shop along with their own customized front end application. 

That being said, a comfortable development workflow was established following microservice principles and 12 Factor app best practices. 

Tools leveraged to create a cohesive and approchable development environment are:
- ![Golang]()
- ![Docker](https://www.docker.com/products/docker-desktop/)
- ![Kubernetes in Docker (KiND)](https://kind.sigs.k8s.io/) 
- ![kustomize (production specifications)](https://kustomize.io/)
- ![Makefile](https://www.gnu.org/software/make/manual/make.html#Makefiles)

## Design

![alt text](https://github.com/tbanta5/CMSC/blob/main/.readme-images/Initia-arch.png)

## Development

### Setup
### Development Patterns

## Educational Model
More than just establishing a development pattern, the hope is that this outline could model potential classroom lab configurations. The pervasiveness, ease of use and accessibility of container technology coupled with the free allowances from ![Docker for Academic institutions](https://www.docker.com/community/open-source/application/) make for a perfect way to distribute software to students with low overhead, little static and expedited usage without the hassle of setup - as is very common in most Java based labs.  

A microservice environment can easily grow and expand with students as they learn new concepts and materials. A single language would be taught throughout the life of the university cirriculum with growing complexities of infrastructure, concurrency and parallelism, web concepts, container technology (fundamentally Unix/Linux OS technologies, security practices and shell scripting), data structures and algorithms, security, compliance and container orchestration systems (Kubernetes) to prepare students for a wholistic "real world" software development workforce. 

Language selection is important. Golang is ubiquitous in container technology and dominates the majority of ![cloud native projects](https://jonathonhenderson.co.uk/2023/07/16/cncf-projects-by-language). It is just as easy to read and write as python with performances of C++ and C. This is why Go is also found in imbedded architectures like Kubernetes or Docker. To point out the obvious, the complexities that come with the JVM, resource heavy environments like Eclipse and Netbeans, and sheer size of the Java language make it less user friendly. Why not substitute a language that is easier to read and write, just as performant and ubiquitous to put students on the leading edge of education? 