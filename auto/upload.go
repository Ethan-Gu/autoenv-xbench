package auto

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/sftp"
)

// Upload file with sftp client
func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	checkErr(err)

	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)

	}

	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	checkErr(err)

	dstFile.Write(ff)
	dstFile.Chmod(0777)
}

// Upload directory with sftp client
func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) {

	localFiles, err := ioutil.ReadDir(localPath)
	checkErr(err)

	//sftpClient.Mkdir(remotePath)
	sftpClient.MkdirAll(remotePath)

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())

		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			uploadDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}
}
