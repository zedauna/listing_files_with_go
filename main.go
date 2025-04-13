/*
* Script de renommage d'extensions de fichiers
*
* Auteur : Jeros VIGAN
* Email :zedauna@programmer.net
* Création : 13/04/2025
* Dernière modification : 13/04/2025
* Version : 1.0.0
*
* Description :
*   Ce script permet de rechercher et renommer massivement
*   les extensions de fichiers de manière interactive.
*
 */
package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func getpath() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Veuillez entrer le chemin du dossier : ")
		dossierPath, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("erreur de lecture : %v", err)
		}
		// Nettoyage du chemin
		dossierPath = strings.TrimSpace(dossierPath)
		dossierPath = strings.TrimRight(dossierPath, "\r\n")     // Gestion des retours chariot Windows
		dossierPath = strings.ReplaceAll(dossierPath, `\`, "\\") // Normalise les séparateurs
		dossierPath = filepath.Clean(dossierPath)

		// Vérification que c'est bien un dossier
		info, err := os.Stat(dossierPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Le dossier '%s' n'existe pas. Veuillez réessayer.\n", dossierPath)
				continue
			}
			return "", fmt.Errorf("erreur d'accès : %v", err)
		}

		if !info.IsDir() {
			fmt.Printf("'%s' n'est pas un dossier valide. Veuillez réessayer.\n", dossierPath)
			continue
		}
		// Conversion en chemin absolu
		absPath, err := filepath.Abs(dossierPath)
		if err != nil {
			return "", fmt.Errorf("erreur de conversion en chemin absolu : %v", err)
		}
		return absPath, nil
	}
}

func listFiles(dir string, extension string) []string {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if len(extension) != 0 {
			if !d.IsDir() && filepath.Ext(path) == extension {
				files = append(files, path)
			}
		} else {
			if !d.IsDir() {
				files = append(files, path)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	return files
}

func listFilesInfo(filePath string) {
	// Nettoyage plus robuste
	filePath = strings.TrimSpace(filePath) // Supprime tous les espaces/retours
	filePath = filepath.Clean(filePath)    // Normalise le chemin
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Erreur : le fichier '%s' n'existe pas.\n", filePath)
		} else {
			fmt.Printf("Erreur d'accès : %v\n", err)
		}
		return
	}
	if fileInfo.IsDir() {
		fmt.Printf("Erreur : '%s' est un dossier, pas un fichier.\n", filePath)
		return
	}

	fmt.Printf("\nInfos fichier :\n")
	fmt.Printf("- Chemin : %s\n", filepath.ToSlash(filePath)) // Standardise les slashs
	fmt.Printf("- Taille : %d octets\n", fileInfo.Size())
	fmt.Printf("- Modifié le : %s\n", fileInfo.ModTime().Format("2006-01-02 15:04:05"))
}

func changerExtension(path string, nouvelleExt string) (string, error) {
	var newPath string
	// Séparation de l'extension
	ext := filepath.Ext(path)
	base := path[:len(path)-len(ext)]

	// Formatage de la nouvelle extension
	if nouvelleExt != "" && !strings.HasPrefix(nouvelleExt, ".") {
		nouvelleExt = "." + nouvelleExt
	}

	// Renommage avec gestion des attributs Windows
	if nouvelleExt != ext && nouvelleExt != "" {
		newPath := base + nouvelleExt
		err := syscall.Rename(path, newPath)
		if err != nil {
			return "", fmt.Errorf("échec du renommage Windows: %v", err)
		}
	} else {
		newPath := base
		err := syscall.Rename(path, newPath)
		if err != nil {
			return "", fmt.Errorf("échec du renommage Windows: %v", err)
		}
	}

	fmt.Printf("\nChangement Extension du Fichier :\n")
	fmt.Printf("- Chemin : %s\n", path)
	fmt.Printf("- Sans extension : %s\n", base)
	fmt.Printf("- Ancienne extension : %s\n", ext)
	fmt.Printf("- Nouvelle extension : %s\n", nouvelleExt)
	// fmt.Printf("- Resultat : %s\n", newPath)
	return newPath, nil
}

func scanner_files(dir string) {
	//Demande de l'extension
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Veuillez entrer l'extension à scanner : ")
	ext_file, _ := reader.ReadString('\n')
	ext_file = strings.TrimSpace(ext_file)
	if !strings.HasPrefix(ext_file, ".") {
		ext_file = "." + ext_file
	}

	if ext_file == "." {
		ext_file = ""
	}
	// Nettoyage plus robuste
	dir = strings.TrimSpace(dir) // Supprime tous les espaces/retours
	dir = filepath.Clean(dir)    // Normalise le chemin

	fmt.Printf("\nParamètre :\n")
	fmt.Printf("- Dossier sélectionné :\n%s\n", dir)
	fmt.Printf("- Extension : %s\n", ext_file)

	files := listFiles(dir, ext_file) //  ".mp4" ,".part"
	for _, v := range files {
		// fmt.Println(v)
		listFilesInfo(v)
		// changerExtension(v, "")
	}
}

func main() {
	var path string
	path, err := getpath()
	if err != nil {
		fmt.Printf("Erreur : %v\n", err)
	}
	//fmt.Printf("\nDossier sélectionné :\n%s\n", path)
	scanner_files(path)
}
