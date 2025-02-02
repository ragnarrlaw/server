set -e

mkdir -p postgres/data

mkdir -p osrm/{data,profiles}

curl -o osrm/profiles/car.lua https://raw.githubusercontent.com/Project-OSRM/osrm-backend/master/profiles/car.lua

wget -P osrm/data http://download.geofabrik.de/asia/sri-lanka-latest.osm.pbf

docker run --rm -v "$(pwd)/osrm/data:/data" ghcr.io/project-osrm/osrm-backend \
    osrm-extract -p /opt/car.lua /data/sri-lanka-latest.osm.pbf

docker run --rm -v "$(pwd)/osrm/data:/data" ghcr.io/project-osrm/osrm-backend \
    osrm-partition /data/sri-lanka-latest.osrm

docker run --rm -v "$(pwd)/osrm/data:/data" ghcr.io/project-osrm/osrm-backend \
    osrm-customize /data/sri-lanka-latest.osrm

docker compose up -d osrm
