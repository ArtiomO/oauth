# Oauth2.0 server.

    Service created as portfolio project so there is no 100% test coverage, simple in memory datastore for clients.
    Swagger page generated also as example.


# Environment variables

    STATIC_URI: /oauth/static/ 
    REDIS_HOST: localhost:6379


# Docker

Server listening on port 8090 inside container.

Build:

    podman build . -t oauth:0

# HOW TO

    You can use swagger to communicate with api http://localhost:5000/swagger/index.html
    

    Launch server container:

        podman create -p 8090:8090 --name oauth --network=local -e STATIC_URI=/static/ -e REDIS_HOST=redis:6379 oauth:0
    
