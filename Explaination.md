# Workflow Possible 

## 1 - Initialisation

Pour le noeud "de départ":
- 1.1 Dans un premier temps, le schéma initial est parsé depuis le(s) fichier(s) qui définissent les schéma. Le schéma est ensuite inscrit en DB. 
  - 1.1.1 Parsage (et validation) schéma
  - 1.1.2 Génération des EntityType et des Attributs en BDD 

## 2 - Pour chaque "hit" sur une api qui modifie un objet

- 2.1 On vérifie que la modification est valide du point de vue du schéma sinon on retourne un erreur. 
  - 2.1.1 Récupérer l'EntityType concerné,
  - 2.1.2 Récupérer les Attributs concernés,
- 2.2.1 Matcher les fields JSON reçus avec les Attributs
- 2.2.2 Préparer et lancer la modification  (?)

## 3 - Pour chaque "hit" sur une api qui lit un objet

- 3.1 Vérifier que l'Entity existe en mémoire
- 3.2 Construire un json (on pourrait [contruire dynamiquement une struct](https://github.com/Ompluscator/dynamic-struct) et le renvoyer)