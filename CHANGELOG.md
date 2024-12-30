# Changelog

## [1.22.0](https://github.com/pmqueiroz/umbra/compare/v1.21.0...v1.22.0) (2024-12-30)


### Features

* add expressions lines in error annotation ([4ca0cae](https://github.com/pmqueiroz/umbra/commit/4ca0cae532315a148298b7c58047e2d13abc748d))
* expressions error annotation ([e654f9f](https://github.com/pmqueiroz/umbra/commit/e654f9f353584859b630232f0f6cec2e3ec3cc40))
* fill get locs ([fb37310](https://github.com/pmqueiroz/umbra/commit/fb373102190b69335f310f7a9bece7914321e6f1))
* statements error annotations ([18102c9](https://github.com/pmqueiroz/umbra/commit/18102c987d6572b3b3df7d3cd5b6541e65a0fcfd))
* use global node as expression and statement ([81a4516](https://github.com/pmqueiroz/umbra/commit/81a4516bff3255552107abf24bd36a33cc047f96))

## [1.21.0](https://github.com/pmqueiroz/umbra/compare/v1.20.1...v1.21.0) (2024-12-13)


### Features

* **math:** add min and max functions ([9ad149b](https://github.com/pmqueiroz/umbra/commit/9ad149bcc1474da57eec6556881c504fd4663cf1))

## [1.20.1](https://github.com/pmqueiroz/umbra/compare/v1.20.0...v1.20.1) (2024-12-13)


### Bug Fixes

* arr to str conversion breaking on more than 1 deep ([0564350](https://github.com/pmqueiroz/umbra/commit/056435073c178cf06957bc40e61eec15fbb1b230))

## [1.20.0](https://github.com/pmqueiroz/umbra/compare/v1.19.0...v1.20.0) (2024-12-12)


### Features

* **arrays:** flatten fns ([918fe26](https://github.com/pmqueiroz/umbra/commit/918fe267702e6ccd175ad02a8c4524ebce9ef40f))
* stack lib ([cd6aaf2](https://github.com/pmqueiroz/umbra/commit/cd6aaf29be945921e334947e7b27b4e6090c66b6))

## [1.19.0](https://github.com/pmqueiroz/umbra/compare/v1.18.0...v1.19.0) (2024-12-12)


### Features

* **arrays:** reverse and clone fns ([9b485cf](https://github.com/pmqueiroz/umbra/commit/9b485cf9989cbcfc232f08a53126ffb8b9275987))
* queue lib ([babc586](https://github.com/pmqueiroz/umbra/commit/babc5860a72595043c7dfc6412bf454e5267e2e5))

## [1.18.0](https://github.com/pmqueiroz/umbra/compare/v1.17.0...v1.18.0) (2024-12-12)


### Features

* is expression ([558b846](https://github.com/pmqueiroz/umbra/commit/558b84669e6cab7df99fd8538cdc7b7829207099))

## [1.17.0](https://github.com/pmqueiroz/umbra/compare/v1.16.1...v1.17.0) (2024-12-12)


### Features

* **set:** remove function ([acd3287](https://github.com/pmqueiroz/umbra/commit/acd328719154bdf4848151642b1143045d1ce825))

## [1.16.1](https://github.com/pmqueiroz/umbra/compare/v1.16.0...v1.16.1) (2024-12-12)


### Bug Fixes

* **set:** null typo ([514c007](https://github.com/pmqueiroz/umbra/commit/514c007fc74ee0c6200014a91361c8de8eb6b26f))

## [1.16.0](https://github.com/pmqueiroz/umbra/compare/v1.15.1...v1.16.0) (2024-12-12)


### Features

* support backwards loop ([8da036d](https://github.com/pmqueiroz/umbra/commit/8da036de39f2484a86ee206a9f643b57f5854b23))


### Performance Improvements

* prevent evaluating step and stop loop every run ([a5c0d85](https://github.com/pmqueiroz/umbra/commit/a5c0d85a5673da17da420a87071d285766482496))
* prevent set loop step if break ([a5c2384](https://github.com/pmqueiroz/umbra/commit/a5c2384dcc62f857a3d8c33236b6073834a2db57))

## [1.15.1](https://github.com/pmqueiroz/umbra/compare/v1.15.0...v1.15.1) (2024-12-12)


### Bug Fixes

* prevent pattern matching return to be supressed ([dcbb03a](https://github.com/pmqueiroz/umbra/commit/dcbb03af19c4979fadd3412e523ff5089537506d))

## [1.15.0](https://github.com/pmqueiroz/umbra/compare/v1.14.1...v1.15.0) (2024-12-11)


### Features

* **arrays:** pass index to map fn ([54c0872](https://github.com/pmqueiroz/umbra/commit/54c08722aa29a8ce0e62c9788200a8149708f72a))
* inline functions ([6d56e34](https://github.com/pmqueiroz/umbra/commit/6d56e34822409dfcc03674ee39e38d4665e4fa22))
* pattern matching case as inline functions ([fa27afb](https://github.com/pmqueiroz/umbra/commit/fa27afb8b8889a98fd6ce55dcd474195e2acc5b3))
* **range:** range over numbers generating arrays of indexes to the count ([43ef7fd](https://github.com/pmqueiroz/umbra/commit/43ef7fd391c5530e66de7296c5166ecf6cb915db))

## [1.14.1](https://github.com/pmqueiroz/umbra/compare/v1.14.0...v1.14.1) (2024-12-11)


### Bug Fixes

* prevent range to return unknown type ([75b1bab](https://github.com/pmqueiroz/umbra/commit/75b1bab16f2ac899e95962b50e50344ceaf3deca))

## [1.14.0](https://github.com/pmqueiroz/umbra/compare/1.13.0...v1.14.0) (2024-12-11)


### Features

* add plus equal and minus equal operators ([07c8a89](https://github.com/pmqueiroz/umbra/commit/07c8a899f0030ba106a0febf13cc3e33883594ce))
* prevent reassigning constant variables ([5e1176b](https://github.com/pmqueiroz/umbra/commit/5e1176b11fd05a244785b39683587fe70b093bc1))
* set lib ([82a9470](https://github.com/pmqueiroz/umbra/commit/82a94702a7e5ba349bb4a6244eafb2bc10501a9a))
