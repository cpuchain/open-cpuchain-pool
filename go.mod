module github.com/sammy007/open-ethereum-pool

go 1.21

replace github.com/ethereum/ethash => github.com/mkrufky/ethereum-ethash v0.0.0-20190805120432-b254c6743dc5

replace github.com/ethereum/go-ethereum => github.com/etclabscore/core-geth v1.12.20

require (
	github.com/ethereum/ethash v0.0.0-20221028165206-dc3eda17d27f
	github.com/ethereum/go-ethereum v1.10.3
	github.com/gorilla/mux v1.8.0
	github.com/yvasiyarov/gorelic v0.0.7
	gopkg.in/redis.v3 v3.6.4
)

require (
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/garyburd/redigo v1.6.4 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/yvasiyarov/go-metrics v0.0.0-20150112132944-c25f46c4b940 // indirect
	github.com/yvasiyarov/newrelic_platform_go v0.0.0-20160601141957-9c099fbc30e9 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/sys v0.16.0 // indirect
	gopkg.in/bsm/ratelimit.v1 v1.0.0-20170922094635-f56db5e73a5e // indirect
)
