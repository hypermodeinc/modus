module github.com/hypermodeinc/modus/runtime

go 1.23.0

require github.com/hypermodeinc/modus/lib/manifest v0.0.0

require github.com/hypermodeinc/modus/lib/metadata v0.0.0

require github.com/hypermodeinc/modus/lib/wasmextractor v0.0.0 // indirect

replace github.com/hypermodeinc/modus/lib/manifest => ../lib/manifest

replace github.com/hypermodeinc/modus/lib/metadata => ../lib/metadata

replace github.com/hypermodeinc/modus/lib/wasmextractor => ../lib/wasmextractor

require (
	github.com/OneOfOne/xxhash v1.2.8
	github.com/archdx/zerolog-sentry v1.8.4
	github.com/aws/aws-sdk-go-v2 v1.32.2
	github.com/aws/aws-sdk-go-v2/config v1.28.0
	github.com/aws/aws-sdk-go-v2/service/s3 v1.66.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.2
	github.com/buger/jsonparser v1.1.1
	github.com/chewxy/math32 v1.11.1
	github.com/dgraph-io/dgo/v230 v230.0.1
	github.com/docker/docker v27.3.1+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/fatih/color v1.18.0
	github.com/getsentry/sentry-go v0.29.1
	github.com/go-viper/mapstructure/v2 v2.2.1
	github.com/goccy/go-json v0.10.3
	github.com/gofrs/flock v0.12.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/renameio v1.0.1
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.1
	github.com/jensneuse/abstractlogger v0.0.4
	github.com/joho/godotenv v1.5.1
	github.com/lestrrat-go/jwx v1.2.30
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/common v0.60.0
	github.com/rs/cors v1.11.1
	github.com/rs/xid v1.6.0
	github.com/rs/zerolog v1.33.0
	github.com/spf13/cast v1.7.0
	github.com/stretchr/testify v1.9.0
	github.com/tetratelabs/wazero v1.8.1
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
	github.com/viterin/vek v0.4.2
	github.com/wundergraph/graphql-go-tools/execution v1.0.7
	github.com/wundergraph/graphql-go-tools/v2 v2.0.0-rc.108
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0
	google.golang.org/grpc v1.67.1
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.41 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.2 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dop251/goja v0.0.0-20240919115326-6c7d1df7ff05 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/pprof v0.0.0-20240925223930-fa3061bff0bc // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jensneuse/byte-template v0.0.0-20231025215717-69252eb3ed56 // indirect
	github.com/kingledion/go-tools v0.6.0 // indirect
	github.com/klauspost/compress v1.17.10 // indirect
	github.com/lestrrat-go/backoff/v2 v2.0.8 // indirect
	github.com/lestrrat-go/blackmagic v1.0.2 // indirect
	github.com/lestrrat-go/httpcc v1.0.1 // indirect
	github.com/lestrrat-go/iter v1.0.2 // indirect
	github.com/lestrrat-go/option v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/phf/go-queue v0.0.0-20170504031614-9abe38d0371d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/r3labs/sse/v2 v2.10.0 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/tidwall/jsonc v0.3.2 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/viterin/partial v1.1.0 // indirect
	github.com/wundergraph/astjson v0.0.0-20240910140849-bb15f94bd362 // indirect
	github.com/wundergraph/cosmo/composition-go v0.0.0-20240926091419-7c3781f4f507 // indirect
	github.com/wundergraph/cosmo/router v0.0.0-20240926091419-7c3781f4f507 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.55.0 // indirect
	go.opentelemetry.io/otel v1.30.0 // indirect
	go.opentelemetry.io/otel/metric v1.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.30.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240924160255-9d4c2d233b61 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240924160255-9d4c2d233b61 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
	rogchap.com/v8go v0.9.0 // indirect
)
