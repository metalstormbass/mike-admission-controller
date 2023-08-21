docker rm -f mike    
docker build . -t michaelbraunbass/mike-admission-controller:main
docker run --network=mike --name=mike -p 8118:8118  -e PORT=8118 -d michaelbraunbass/mike-admission-controller:main