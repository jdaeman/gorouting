# Go simple routing engine 

## Quick start
1. Prepare osm pbf data
    - https://download.geofabrik.de/asia/south-korea.html
1. Extract pbf data
1. Run backend
1. Run frontend

--
## docker build
1. docker build . -f docker/Dockerfile -t gorouting:{version}
1. docker run ...

---
## Extractor
1. install pkg-config for window
    - ;;
1. 

---
## Routed


---
## Frontend
1. https://github.com/Project-OSRM/osrm-frontend
1. Fix serveral configuration.
    - src/leaflet_options.js
        - services
            - path => ip:port
    - docker/Dockerfile
        - ENV OSRM_BACKEND => ip:port
1. issues
    - container 실행 시, SyntaxError: Unexpected identifier 
    - 해결 방법.
    - https://github.com/Project-OSRM/osrm-frontend/issues/350

1. --env
    - container 실행 시, 환경 변수를 통한 backend 설정
    - --env OSRM_BACKEND=http://ip:port

---
## Docker
1. docker-compose.yml
- https://stackoverflow.com/questions/65285379/docker-volume-mapping-windows-incredible-slow
- https://stackoverflow.com/questions/36249744/interactive-shell-using-docker-compose