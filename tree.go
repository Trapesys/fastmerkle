package fastmerkle

// Node represents a single node in the Merkle
type Node struct {
	Hash []byte // The hash of the children nodes

	Left   *Node // Reference to the left child
	Right  *Node // Reference to the right child
	Parent *Node // Reference to the parent node
}

// duplicate creates a copy of the node
// and its reference peers
func (n *Node) duplicate() *Node {
	return &Node{
		Hash:   n.Hash,
		Left:   n.Left,
		Right:  n.Right,
		Parent: n.Parent,
	}
}

// GetHash returns the hash of the node's children
func (n *Node) GetHash() []byte {
	return n.Hash
}

// MerkleTree represents the perfect Merkle binary tree
type MerkleTree struct {
	Root *Node // The root of the Merkle binary tree
}

// GetRoot returns the root of the Merkle tree
func (m *MerkleTree) GetRoot() *Node {
	return m.Root
}

// GetRootHash returns the root hash fo the Merkle tree
func (m *MerkleTree) GetRootHash() []byte {
	return m.Root.GetHash()
}
