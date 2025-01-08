module github.com/hypermodeinc/modus/runtime

go 1.23.1

toolchain go1.23.4

require (
	github.com/hypermodeinc/modus/lib/manifest v0.16.0
	github.com/hypermodeinc/modus/lib/metadata v0.15.0
	github.com/hypermodeinc/modus/lib/wasmextractor v0.13.0 // indirect
)

require (
	github.com/OneOfOne/xxhash v1.2.8
	github.com/archdx/zerolog-sentry v1.8.5
	github.com/aws/aws-sdk-go-v2 v1.32.7
	github.com/aws/aws-sdk-go-v2/config v1.28.7
	github.com/aws/aws-sdk-go-v2/service/s3 v1.72.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.3
	github.com/buger/jsonparser v1.1.1
	github.com/chewxy/math32 v1.11.1
	github.com/dgraph-io/dgo/v240 v240.1.0
	github.com/docker/docker v27.4.1+incompatible
	github.com/docker/go-connections v0.5.0
	github.com/fatih/color v1.18.0
	github.com/getsentry/sentry-go v0.31.1
	github.com/go-viper/mapstructure/v2 v2.2.1
	github.com/goccy/go-json v0.10.4
	github.com/gofrs/flock v0.12.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/google/renameio/v2 v2.0.0
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/jensneuse/abstractlogger v0.0.4
	github.com/joho/godotenv v1.5.1
	github.com/lestrrat-go/jwx v1.2.30
	github.com/neo4j/neo4j-go-driver/v5 v5.27.0
	github.com/prometheus/client_golang v1.20.5
	github.com/prometheus/common v0.61.0
	github.com/puzpuzpuz/xsync/v3 v3.4.0
	github.com/rs/cors v1.11.1
	github.com/rs/xid v1.6.0
	github.com/rs/zerolog v1.33.0
	github.com/spf13/cast v1.7.1
	github.com/stretchr/testify v1.10.0
	github.com/tetratelabs/wazero v1.8.2
	github.com/tidwall/gjson v1.18.0
	github.com/tidwall/sjson v1.2.5
	github.com/viterin/vek v0.4.2
	github.com/wundergraph/graphql-go-tools/execution v1.2.0
	github.com/wundergraph/graphql-go-tools/v2 v2.0.0-rc.137
	golang.org/x/exp v0.0.0-20241217172543-b2144cdd0a67
	golang.org/x/sys v0.29.0
	google.golang.org/grpc v1.69.2
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.7 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.48 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.26 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.26 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.7 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/containerd/log v0.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.3.0 // indirect
	github.com/distribution/reference v0.6.0 // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dop251/goja v0.0.0-20241024094426-79f3a7efcdbd // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-sourcemap/sourcemap v2.1.4+incompatible // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jensneuse/byte-template v0.0.0-20231025215717-69252eb3ed56 // indirect
	github.com/kingledion/go-tools v0.6.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
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
	github.com/wundergraph/astjson v0.0.0-20241210135722-15ca0ac078f8 // indirect
	github.com/wundergraph/cosmo/composition-go v0.0.0-20241223134725-4acccc1dcaca // indirect
	github.com/wundergraph/cosmo/router v0.0.0-20241223134725-4acccc1dcaca // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.58.0 // indirect
	go.opentelemetry.io/otel v1.33.0 // indirect
	go.opentelemetry.io/otel/metric v1.33.0 // indirect
	go.opentelemetry.io/otel/trace v1.33.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20241223144023-3abc09e42ca8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241223144023-3abc09e42ca8 // indirect
	google.golang.org/protobuf v1.36.1 // indirect
	gopkg.in/cenkalti/backoff.v1 v1.1.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/v3 v3.5.1 // indirect
	rogchap.com/v8go v0.9.0 // indirect
)
