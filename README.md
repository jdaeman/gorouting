# Go simple routing engine 

## Quick start
1. Prepare osm pbf data
    - https://download.geofabrik.de/asia/south-korea.html
1. Extract pbf data
1. Run backend
1. Run frontend

---
## Extractor
1. install pkg-config for window
    - ;;
1. uu

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
    - https://github.com/Project-OSRM/osrm-frontend/issues/350

---
## Docker
1. docker-compose.yml
