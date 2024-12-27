# go-organizer

https://github.com/user-attachments/assets/95ea4ec8-2111-443a-8c5a-c8d6caa9f907


![Screenshot 2024-12-27 155148](https://github.com/user-attachments/assets/1a617bf5-6e7b-456b-895b-b7e1eb27a0c3)


### Version "Overlay"

![Screenshot 2024-12-27 155155](https://github.com/user-attachments/assets/06567361-8fba-4be0-8e25-81aeea7c00c6)

### Little to no memory usage

![Screenshot 2024-12-27 160148](https://github.com/user-attachments/assets/877f46a3-d774-4a4d-83b8-2c710ac7c739)

## FR (English below)

### Un organiseur pour Dofus Unity en Golang, **CECI N'EST PAS UN LOGICIEL DE TRICHE**, aucune automatisation de clic / déplacements / échanges etc..
### Sa seule fonction est de permettre l'organisation des fenêtres des clients Dofus, sauvegarder un ordre et sauvegarder trois raccourcis.

- "Toggle Organizer": Active/Désactive l'organiseur
- "Previous Character": Active la fenêtre du personnage suivant
- "Next Character": Active la fenêtre du personnage précédent

Les boutons :
- "Pin to Top": Épingle la fenêtre au premier plan.
- "Fetch": Lance une recherche de fenêtre Dofus, par exemple : vous avez lancé un nouveau compte, utilisez ce bouton pour rafraîchir la liste.
- "Load": Organise les personnages en suivant l'ordre sauvegardé via le bouton "Save", si un personnage n'est pas dans la liste sauvegardée, il sera alors placé à la fin de celle-ci.
- "Save": Sauvegarde l'ordre de priorité, si 8 comptes sont sauvegardés, mais seulement 4 sont connectés, l'ordre devrait être préservé en se basant sur la liste des 8.
- "Switch to Overlay Mode" : Passe l'application au premier plan et réduit au maximum sa taille, un espace blanc sur la droite est disponible pour la déplacer.

Autre Fonctionnalité :
- Cliquer sur le nom d'un personnage passe alors sa fenêtre au premier plan
- L'application s'arrête automatiquement au bout d'une heure sans utilisation.

## EN

### An organizer for Dofus Unity in Golang, **THIS IS NOT A CHEATING SOFTWARE**, no clicking automation / no auto-move or auto-follow / no auto-trade etc..
### Its only use is to facilitate organizing Dofus Windows, save them in order and also bind 3 keybinds to actions.

- "Toggle Organizer": Toggle On/Off the Organizer
- "Previous Character": Activate Previous Character Window
- "Next Character": Activate Next Character Window

The buttons:
- "Pin to Top": Pin go-organizer window to always be on top.
- "Fetch": Attempt to Fetch all Dofus Windows, ex: if you logged a new character, use this button to update the list.
- "Load": Will Order the windows based on the saved order if it exists, if a character is not currently in the saved list, he will be placed at the end of the list.
- "Save": Save the current set order, if 8 characters are saved, but only 4 are currently logged in, the order will be preserved based on the 8 saved characters.
- "Switch to Overlay Mode" : Set the app to be always on top and reduce its size to the bare minimum, has a small blank space on the right side which can be used to drag it around.

Other functionality:
- Clicking on a character name will activate its window to the foreground.
- The app will stop listening to keybinds after one hour of inactivity.


Knowns Bugs :
- Mouse buttons seems to have a hard time triggering the window change

