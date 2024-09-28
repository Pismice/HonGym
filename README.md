# Pages
### Misc
## Login/Inscription
Si l utilisteur n est pas connecte.

### Activite
## Accueil (workout en cours)
On voit le workout en cours ainsi que les differences seances effectuees et celles qui restent a faire.
En cliquant sur une seance pas effectuee on peut la lancer.
On affiche egalement combien de semaines il reste avant la fin de ce workout programme.

## Seance en cours
Même principe que pour accueil mais avec les exercices de la seance en cours.

## Exercice en cours
"Serie X sur Y"
On rentre le nombre de rep ainsi que le poids utilise puis on passe a la serie suivante jusqu a avoir accompli les Y series puis on retourne sur *Seance en cours*
On affiche egalement le PR pour cet exercice, il faut faire 8 reps pour valider un nouveau PR.

### Gestion des entites
## Exercices
On peut voir tous les exercices actuels de l utilisateur (possiblite de modif) puis un bouton permet d en creer de nouveaux.

## Seances
Meme principe que pour *Exercices* mais on peut rajouter ou enelver des exercices.
Possiblite de superset a garder en tete

## Workouts
Meme principe que pour *Seances* mais on peut rajouter ou enelver des seances ainsi que donner un nombre de semaines a effectuer pour ce workout.

# Implementation
Gerer la session utilistaure avec un cookie et middleware. (meme maniere que Zig conquest)
Chaque page correspond a un /GET
