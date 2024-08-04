package Uploads

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"shopperia/src/common/models"
	UserDTO "shopperia/src/core/user/domain/DTO"
)

func (US *UploadService) makeNewCollection(userRepositoryPath, collectionName string) (string, error) {
	if userRepositoryPath == "" || collectionName == "" {
		return "", errors.New("collection name or collection name is empty")
	}
	collectionPath, collErr := US.MakeNewDirectory(userRepositoryPath, collectionName)
	if collErr != nil {
		return "", collErr
	}

	return collectionPath, nil
}

func (US *UploadService) CreateCollection(form UserDTO.CreateCollectionForm) (models.CollectionData, error) {
	if form.UserRepositoryPath == "" || form.CollectionName == "" {
		errStr := fmt.Sprint("Please provide all fields")
		return models.CollectionData{}, errors.New(errStr)
	}

	collection, err := US.makeNewCollection(form.UserRepositoryPath, form.CollectionName)
	if err != nil {
		return models.CollectionData{}, err
	}

	newCollectionId := uuid.New()

	collectionData := models.CollectionData{
		CollectionId:   newCollectionId,
		UserRepository: form.UserRepositoryPath,
		CollectionName: form.CollectionName,
		CollectionPath: collection,
	}

	return collectionData, nil
}

func (US *UploadService) InsertImageOnCollection(repositoryPath, collectionPath string, image models.UploadImageForm) (models.ImageData, error) {
	if collectionPath == "" {
		return models.ImageData{}, errors.New("collection path is empty")
	}

	requestChan := make(chan *uploadImageSingleAttempt)

	go US.uploadWorker(requestChan)

	Path := filepath.Join(US.OriginPath, repositoryPath, collectionPath)

	uploadAttempt := &uploadImageSingleAttempt{
		Image:         image,
		Done:          make(chan struct{}),
		DirectoryPath: Path,
	}

	requestChan <- uploadAttempt

	<-uploadAttempt.Done
	if uploadAttempt.Status == nil {
		return models.ImageData{}, uploadAttempt.Status
	}

	imageData := uploadAttempt.Data

	return imageData, nil
}

type uploadImageAttempt struct {
	UserId         uuid.UUID
	RepositoryPath string
	CollectionPath string
	FileName       string
	FileExtension  string
	Image          bytes.Buffer
}

type resolveImageAttempt struct {
	ImageData models.ImageData
	Status    error
}

func (US *UploadService) uploadImage(attemptRequest uploadImageAttempt, resolveChan chan<- resolveImageAttempt) {
	if attemptRequest.FileName == "" || attemptRequest.FileExtension == "" || attemptRequest.UserId == uuid.Nil || attemptRequest.RepositoryPath == "" || attemptRequest.CollectionPath == "" {
		resolveChan <- resolveImageAttempt{
			Status: errors.New("paremeters required"),
		}
	}

	collectionPath := filepath.Join(attemptRequest.CollectionPath)
	imageData, err := US.upload(collectionPath, attemptRequest.FileName, attemptRequest.FileExtension, attemptRequest.Image, attemptRequest.UserId)
	resolve := resolveImageAttempt{
		ImageData: imageData,
		Status:    err,
	}

	resolveChan <- resolve
}

func (US *UploadService) reedResolvesChannel(numAttempts int, resolveChan <-chan resolveImageAttempt) ([]models.ImageData, error) {
	if numAttempts == 0 {
		return nil, errors.New("numAttempts must be 1")
	}

	var imagesData []models.ImageData
	i := 0
	for {
		if i >= numAttempts {
			return imagesData, nil
		}
		select {
		case res := <-resolveChan:
			if res.Status != nil {
				return nil, res.Status
			}

			imagesData = append(imagesData, res.ImageData)
			i++

		}
	}
}

func (US *UploadService) InsertMultipleImagesOnCollection(collectionPath string, forms []models.UploadImageForm) ([]models.ImageData, error) {
	if collectionPath == "" || len(forms) == 0 {
		return nil, errors.New("No parameters provide")
	}

	for i, image := range forms {
		nextImage := forms[i+1]

		if image.FileName == nextImage.FileName {
			nextImage.FileName += "_" + string(i+1)
		}
	}

	resolvesChan := make(chan resolveImageAttempt, len(forms))
	for _, form := range forms {

		uploadAttempt := uploadImageAttempt{
			UserId:         form.UserID,
			CollectionPath: collectionPath,
			FileName:       form.FileName,
			FileExtension:  form.FileExtension,
			Image:          form.ImageData,
		}

		go US.uploadImage(uploadAttempt, resolvesChan)

	}

	Data, err := US.reedResolvesChannel(len(forms), resolvesChan)
	if err != nil {
		return Data, err
	}

	return Data, nil
}

type GetImageattempt struct {
	buf  bytes.Buffer
	err  error
	done chan struct{}
}

func readimageBuffer(numAttempts int, attempchan <-chan GetImageattempt) ([]GetImageattempt, error) {

	var resolves = []GetImageattempt{}

	i := 0

	for {
		if i == numAttempts {
			break
		}

		select {
		case res := <-attempchan:
			if res.err != nil {
				resolves = append(resolves, GetImageattempt{err: res.err})
			}

			resolves = append(resolves, GetImageattempt{buf: res.buf})
			i++
		}

	}

	return resolves, nil
}

func searchImage(imgChan chan<- GetImageattempt, fileName string) {
	if fileName == "" {
		imgChan <- GetImageattempt{
			err: errors.New("file name is empty"),
		}
		return
	}

	var buf bytes.Buffer
	imgsBytes, err := os.ReadFile(fileName)
	if err != nil {
		imgChan <- GetImageattempt{
			err: err,
		}
		return
	}

	_, err = buf.Write(imgsBytes)
	if err != nil {
		imgChan <- GetImageattempt{err: err}
		return
	}

	imgChan <- GetImageattempt{
		buf: buf,
		err: nil,
	}

}

func (US *UploadService) GetAllMediaFromCollection(repositoryPath, collectionPath string, forms []models.GetImageForm) ([]models.GetImage, error) {
	if repositoryPath == "" || collectionPath == "" {
		return nil, errors.New("repository path or collection path is empty")
	}

	requestImageAttempt := make(chan GetImageattempt, len(forms))

	for _, file := range forms {

		if file.FileName == "" || file.FileExtension == "" {
			return nil, errors.New("file name or file extension is empty")
		}

		searchFileName := US.OriginPath + "/" + repositoryPath + "/" + collectionPath + "/" + file.FileName + "." + file.FileExtension

		go searchImage(requestImageAttempt, searchFileName)
	}

	GetImageAttempts, err := readimageBuffer(len(forms), requestImageAttempt)
	if err != nil {
		return nil, err
	}

	var images = []models.GetImage{}

	for i, imageBuf := range GetImageAttempts {
		image := models.GetImage{
			FileName:      forms[i].FileName,
			FileExtension: forms[i].FileExtension,
			ImageBuffer:   imageBuf.buf,
		}

		images = append(images, image)
	}

	close(requestImageAttempt)

	return images, nil
}

func (US *UploadService) DeleteImageOnCollection(request models.DeleteOnCollectionRequest) error {
	if request.UserRepositoryPath == "" || request.CollectionName == "" || request.FileName == "" || request.FileExtension == "" {
		return errors.New("")
	}

	requestChan := make(chan *deleteRequest)

	go US.deleteWorker(1, requestChan)

	fileName := request.FileName + "." + request.FileExtension
	resourcePath := filepath.Join(US.OriginPath, request.UserRepositoryPath, request.CollectionName, fileName)

	req := &deleteRequest{
		ResourcePath: resourcePath,
		Done:         make(chan struct{}),
	}

	requestChan <- req

	<-req.Done
	if req.Status != nil {
		return req.Status
	}

	close(requestChan)
	close(req.Done)

	return nil
}

func (US *UploadService) DeleteMultipleImagesOnCollection(requests []models.DeleteOnCollectionRequest) error {

	requestsChan := make(chan *deleteRequest, len(requests))

	go US.deleteWorker(len(requests), requestsChan)

	for _, request := range requests {
		fileName := request.FileName + "." + request.FileExtension
		resourcePath := filepath.Join(US.OriginPath, request.UserRepositoryPath, request.CollectionName, fileName)

		req := &deleteRequest{
			ResourcePath: resourcePath,
			Done:         make(chan struct{}),
		}

		requestsChan <- req
		<-req.Done
		if req.Status != nil {
			return req.Status
		}
	}

	return nil
}

func (US *UploadService) CheckIfImageIsOnCollection(collectionPath, fileName, fileExtension string) bool {
	if collectionPath == "" || fileName == "" || fileExtension == "" {
		return false
	}

	completeFileName := fileName + "." + fileExtension
	path := filepath.Join(collectionPath, completeFileName)
	_, err := os.Stat(path)
	if !os.IsExist(err) {
		return false
	}

	return true
}

func (US *UploadService) UpdateImageOnCollection(collectionPath, fileName, fileExtension string, form models.UploadImageForm) (models.ImageData, error) {
	if collectionPath == "" || fileName == "" || fileExtension == "" {
		return models.ImageData{}, errors.New("no parameters provide")
	}

	deleteRequestChan := make(chan *deleteRequest)
	updateRequestChan := make(chan *uploadImageSingleAttempt)

	go US.deleteWorker(1, deleteRequestChan)

	completeImagePath := collectionPath + "/" + fileName + "." + fileExtension

	deleteReq := &deleteRequest{
		ResourcePath: completeImagePath,
		Done:         make(chan struct{}),
	}

	deleteRequestChan <- deleteReq

	<-deleteReq.Done
	if deleteReq.Status != nil {
		return models.ImageData{}, deleteReq.Status
	}

	updateAttempt := &uploadImageSingleAttempt{
		Image:         form,
		DirectoryPath: "",
		Status:        nil,
		Done:          nil,
	}

	go US.uploadWorker(updateRequestChan)

	updateRequestChan <- updateAttempt

	<-updateAttempt.Done
	if updateAttempt.Status != nil {
		return updateAttempt.Data, updateAttempt.Status
	}

	return updateAttempt.Data, nil

}

func (US *UploadService) DeleteCollection(userRepository, collectionName string) error {
	return nil
}

func (US *UploadService) updateCollectionWorker(numAttempts int, attempt chan *attemptchangeCollectionName) {
	i := 0
	for {
		if i == numAttempts {
			break
		}
		select {
		case req := <-attempt:

			US.getFileNamesOfCollection(req)
			US.getFilesData(req)
			newCollectionName, err := US.makeNewCollection(req.OldPath, req.NewCollectionName)
			if err != nil {
				req.Status = err
			}

			req.NewCollectionPath = newCollectionName

			US.InsertFileDataOnNewCollection(req)

			err = US.deleteDirectory(req.NewPath)
			if err != nil {
				req.Status = err
			}

			req.Done <- struct{}{}

		}

	}

}

func (US *UploadService) getFileNamesOfCollection(attempt *attemptchangeCollectionName) {

	var fileNames = []string{}

	attempt.Status = filepath.WalkDir(attempt.OldPath, func(path string, d os.DirEntry, err error) error {

		if err != nil {
			attempt.Status = err
		}
		if !d.IsDir() {
			fileNames = append(fileNames, filepath.Base(path))
		}

		return err
	})

	var filesOnCollection = []file{}

	for _, name := range fileNames {
		fileName := name[:len(name)-len(filepath.Ext(name))]
		fileExtension := filepath.Ext(fileName)

		fileOnCollection := file{
			Name:         fileName,
			Extension:    fileExtension,
			CompleteName: fileName + fileExtension,
		}

		filesOnCollection = append(filesOnCollection, fileOnCollection)
	}

	attempt.Files = filesOnCollection

}

func (US *UploadService) getFilesData(attempt *attemptchangeCollectionName) {

	var files []bytes.Buffer

	for _, fileName := range attempt.Files {

		var buf bytes.Buffer

		path := filepath.Join(attempt.OldPath, fileName.CompleteName)
		image, err := os.ReadFile(path)

		if err != nil {
			attempt.Status = err
		}

		buf.Write(image)

		files = append(files, buf)

	}

	attempt.FilesData = files

}

func (US *UploadService) InsertFileDataOnNewCollection(attempt *attemptchangeCollectionName) {
	updatedImages := []models.ImageData{}
	for i := range attempt.Files {

		updatedImage, err := US.upload(attempt.NewPath, attempt.Files[i].Name, attempt.Files[i].Extension, attempt.FilesData[i], attempt.UserID)
		if err != nil {
			attempt.Status = err
		}
		updatedImages = append(updatedImages, updatedImage)
	}

	attempt.ImagesUpdated = updatedImages
}

func (US *UploadService) UpdateCollectionName(userRepository, collectionName, newCollectionName string, userID uuid.UUID) (string, string, error) {
	if userRepository == "" || collectionName == "" || newCollectionName == "" {
		return "", "", errors.New("no parameters provide")
	}

	oldPath := filepath.Join(US.OriginPath, userRepository, collectionName)
	newPath := filepath.Join(US.OriginPath, userRepository, newCollectionName)

	attemptChan := make(chan *attemptchangeCollectionName)

	go US.updateCollectionWorker(1, attemptChan)

	attempt := &attemptchangeCollectionName{
		UserID:            userID,
		OldPath:           oldPath,
		NewPath:           newPath,
		NewCollectionName: newCollectionName,
		Done:              make(chan struct{}),
	}

	attemptChan <- attempt

	<-attempt.Done
	close(attemptChan)
	close(attempt.Done)
	if attempt.Status != nil {
		return "", "", attempt.Status
	}

	return attempt.NewPath, attempt.NewCollectionName, nil
}
