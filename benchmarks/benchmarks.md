|Package|Time|Time %|Allocations|
|-------|----|------|-----------|
|eltorocorp/logger.Info-4 |309  ns/op|    0%|0  allocs/op|
|eltorocorp/logger.InfoD-4         |324  ns/op|    5%|0  allocs/op|
|rs/zerolog.Check-4       |355  ns/op|   15%|0  allocs/op|
|rs/zerolog-4             |357  ns/op|   16%|0  allocs/op|
|Zap-4         |476  ns/op|   54%|0  allocs/op|
|Zap.Check-4   |505  ns/op|   63%|0  allocs/op|
|Zap.Sugar-4   |971  ns/op|  214%|2  allocs/op|
|go-kit/kit/log-4         |3189  ns/op|  932%|24  allocs/op|
|apex/log-4    |8199  ns/op| 2553%|25  allocs/op|
|sirupsen/logrus-4        |8556  ns/op| 2669%|37  allocs/op|
|inconshreveable/log15-4  |9620  ns/op| 3013%|31  allocs/op|
