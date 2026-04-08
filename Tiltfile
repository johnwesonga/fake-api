# Use Kustomize
k8s_yaml('tmp.yaml')

# Build Docker image
docker_build(
    'johnwesonga/fake-api',
    '.',
    live_update=[
        sync('.', '/app'),
        run('go build -o server', trigger=['*.go']),
    ],
)

# Tell Tilt which resource to track
k8s_resource(
    'fake-api',
    port_forwards=9000,
)