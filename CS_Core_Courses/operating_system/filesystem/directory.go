package filesystem

import (
	"fmt"
	"strings"
)

// ============================================================
// 目录管理
// 408考点：目录结构（单级/两级/树形）、路径解析
// ============================================================

// DirEntry 目录项
type DirEntry struct {
	Name     string
	InodeNum int
	IsDir    bool
}

// DirectoryNode 目录树节点
type DirectoryNode struct {
	Name     string
	InodeNum int
	IsDir    bool
	Entries  []*DirectoryNode // 子目录/文件
	Parent   *DirectoryNode
}

// FileSystem 简单文件系统（树形目录结构）
type FileSystem struct {
	Root      *DirectoryNode
	NextInode int
}

// NewFileSystem 创建文件系统
func NewFileSystem() *FileSystem {
	root := &DirectoryNode{
		Name:     "/",
		InodeNum: 0,
		IsDir:    true,
		Entries:  make([]*DirectoryNode, 0),
	}
	return &FileSystem{Root: root, NextInode: 1}
}

// parsePath 解析路径为各级目录名
func parsePath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// findNode 从指定节点开始查找路径
func (fs *FileSystem) findNode(start *DirectoryNode, parts []string) *DirectoryNode {
	current := start
	for _, part := range parts {
		if part == ".." {
			if current.Parent != nil {
				current = current.Parent
			}
			continue
		}
		if part == "." {
			continue
		}
		found := false
		for _, entry := range current.Entries {
			if entry.Name == part {
				current = entry
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return current
}

// ResolvePath 路径解析（绝对路径）
func (fs *FileSystem) ResolvePath(path string) *DirectoryNode {
	if path == "/" {
		return fs.Root
	}
	parts := parsePath(path)
	return fs.findNode(fs.Root, parts)
}

// CreateDir 创建目录
func (fs *FileSystem) CreateDir(path string) bool {
	parts := parsePath(path)
	if len(parts) == 0 {
		return false
	}

	// 找到父目录
	parent := fs.Root
	if len(parts) > 1 {
		parent = fs.findNode(fs.Root, parts[:len(parts)-1])
	}
	if parent == nil || !parent.IsDir {
		return false
	}

	// 检查是否已存在
	name := parts[len(parts)-1]
	for _, entry := range parent.Entries {
		if entry.Name == name {
			return false
		}
	}

	// 创建新目录
	newDir := &DirectoryNode{
		Name:     name,
		InodeNum: fs.NextInode,
		IsDir:    true,
		Entries:  make([]*DirectoryNode, 0),
		Parent:   parent,
	}
	fs.NextInode++
	parent.Entries = append(parent.Entries, newDir)
	return true
}

// CreateFile 创建文件
func (fs *FileSystem) CreateFile(path string) bool {
	parts := parsePath(path)
	if len(parts) == 0 {
		return false
	}

	parent := fs.Root
	if len(parts) > 1 {
		parent = fs.findNode(fs.Root, parts[:len(parts)-1])
	}
	if parent == nil || !parent.IsDir {
		return false
	}

	name := parts[len(parts)-1]
	for _, entry := range parent.Entries {
		if entry.Name == name {
			return false
		}
	}

	newFile := &DirectoryNode{
		Name:     name,
		InodeNum: fs.NextInode,
		IsDir:    false,
		Parent:   parent,
	}
	fs.NextInode++
	parent.Entries = append(parent.Entries, newFile)
	return true
}

// ListDir 列出目录内容
func (fs *FileSystem) ListDir(path string) []DirEntry {
	node := fs.ResolvePath(path)
	if node == nil || !node.IsDir {
		return nil
	}

	entries := make([]DirEntry, 0)
	for _, e := range node.Entries {
		entries = append(entries, DirEntry{
			Name:     e.Name,
			InodeNum: e.InodeNum,
			IsDir:    e.IsDir,
		})
	}
	return entries
}

// PrintTree 打印目录树
func (fs *FileSystem) PrintTree() {
	printTreeHelper(fs.Root, "", true)
}

func printTreeHelper(node *DirectoryNode, prefix string, isLast bool) {
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	typeStr := ""
	if node.IsDir {
		typeStr = "/" // 目录标识
	}

	if node.Parent == nil {
		fmt.Printf("%s (inode=%d)\n", node.Name, node.InodeNum)
	} else {
		fmt.Printf("%s%s%s (inode=%d)\n", prefix, connector, node.Name+typeStr, node.InodeNum)
	}

	childPrefix := prefix
	if isLast {
		childPrefix += "    "
	} else {
		childPrefix += "│   "
	}

	for i, child := range node.Entries {
		printTreeHelper(child, childPrefix, i == len(node.Entries)-1)
	}
}

// DirectoryExample 目录管理示例
func DirectoryExample() {
	fmt.Println("\n--- 目录管理 ---")

	fs := NewFileSystem()

	// 创建目录结构
	fs.CreateDir("home")
	fs.CreateDir("home/user")
	fs.CreateDir("home/user/documents")
	fs.CreateDir("home/user/downloads")
	fs.CreateDir("etc")
	fs.CreateDir("var")
	fs.CreateDir("var/log")

	// 创建文件
	fs.CreateFile("home/user/documents/report.txt")
	fs.CreateFile("home/user/documents/notes.md")
	fs.CreateFile("home/user/downloads/image.png")
	fs.CreateFile("etc/config.ini")
	fs.CreateFile("var/log/system.log")

	fmt.Println("目录树:")
	fs.PrintTree()

	// 路径解析
	fmt.Println("\n路径解析:")
	paths := []string{"/", "/home/user", "/home/user/documents", "/var/log/system.log"}
	for _, p := range paths {
		node := fs.ResolvePath(p)
		if node != nil {
			typeStr := "文件"
			if node.IsDir {
				typeStr = "目录"
			}
			fmt.Printf("  \"%s\" → inode=%d (%s)\n", p, node.InodeNum, typeStr)
		} else {
			fmt.Printf("  \"%s\" → 未找到\n", p)
		}
	}

	// 列出目录内容
	fmt.Println("\nls /home/user/documents:")
	entries := fs.ListDir("/home/user/documents")
	for _, e := range entries {
		typeStr := "文件"
		if e.IsDir {
			typeStr = "目录"
		}
		fmt.Printf("  %-20s inode=%-4d %s\n", e.Name, e.InodeNum, typeStr)
	}
}
