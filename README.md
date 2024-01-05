## Description

This is an example backend twitter-clone implementation with Clean Architecture in Go (Golang).

Rule of Clean Architecture by Uncle Bob

- Independent of Frameworks. The architecture does not depend on the existence of some library of feature laden software. This allows you to use such frameworks as tools, rather than having to cram your system into their limited constraints.
- Testable. The business rules can be tested without the UI, Database, Web Server, or any other external element.
- Independent of UI. The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a console UI, for example, without changing the business rules.
- Independent of Database. You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules are not bound to the database.
- Independent of any external agency. In fact your business rules simply don’t know anything at all about the outside world.

More at https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html

This project has 4 Domain layer :

- Models Layer
- Repository Layer
- Usecase Layer
- Delivery Layer

#### The diagram:

![golang clean architecture](https://github.com/alifahsanilsatria/twitter-clone/raw/master/clean-arch.png)

It may different already, but the concept still the same in application level.

Here are our repository structure 

#### How to run repository
1. Install & run docker engine or desktop in your computer
2. Since this project use docker compose, you can simply run this command:
```docker-compose up --build```
3. If this is the first time you run this docker compose on your computer. skip this step if you run this for second or above times:
a. login to postgresql with the following credentials: 
    * username = twitter_clone
    * password = twitterclone123
    * host = localhost
    * port = 5432

    b. execute every query in twitter.sql
3. If you want to stop & remove containers & networks, run this command:
```docker-compose down```
