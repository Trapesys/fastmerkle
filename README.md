## Overview

`fastmerkle` is fast tooling for generating Merkle trees, as well as generating and verifying Merkle proofs. Currently,
the supported hashing strategy is _Keccak256_. The tree generation handles uneven data sets by duplicating the last
element on the tree level.

![GitHub](https://img.shields.io/github/license/Trapesys/fastmerkle)
![GitHub branch checks state](https://img.shields.io/github/checks-status/Trapesys/fastmerkle/main)
![GitHub all releases](https://img.shields.io/github/downloads/Trapesys/fastmerkle/total)

## Quick Start 📝

The overall API footprint of the package is relatively small.

### Package Install

To install the `fastmerkle` package, run:

```bash
go get github.com/Trapesys/fastmerkle
````

### Generate Merkle Tree

```go
// Arbitrary random data used to construct the tree
var randomData [][]byte
randomData = ...

// Generate the Merkle tree
merkleTree, createErr := GenerateMerkleTree(randomData)
if createErr != nil {
// Error occurred during Merkle tree generation
...
}

// Get the Merkle root hash
var merkleRootHash []byte
merkleRootHash = merkleTree.GetRootHash()
```

---

Copyright 2022 Trapesys

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the
License. You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "
AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
language governing permissions and limitations under the License.
