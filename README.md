# go-organizer

![Screenshot 2024-11-24 141639](https://github.com/user-attachments/assets/0a6eb79a-9ae0-40e4-9bb6-7c0af7995229)

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

Autre Fonctionnalité :
- Cliquer sur le nom d'un personnage passe alors sa fenêtre au premier plan

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

Other functionality:
- Clicking on a character name will activate its window to the foreground.

