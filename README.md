# PreLoRUGo

PreLoRUGO est le bakc-end écrit en go du logiciel PreLoRU. Il est utilisé en complément avec prelorufront son front end écrit en Javascript avec vue.js et vuetify.

PreloRUGO respecte globalement les principes REST pour son API. Cependant, elles ne respectent pas la totalité des normes REST. Il utilise aussi le fait que le back end et le front end sont développés en même temps et que certaines requêtes ont été conçues de manière intégrée. Par exemple, pour les mises à jour, l'ID n'est pas forcément utilisé dans le chemin de routage de la requête de mise à jour.

La sécurité est assurée par des token JWT qui sont générés par le back-end. Les fonctions de gestion des délais de validité des tokens et de régénération sont implémentées directement dans le back-end. Un minimum de permanence des tokens est assuré par un fichier json qui est sauvegardé lorsque le serveur s'arrête afin de permettre de conserver les connexions actives y compris lors d'une mise à jour du serveur.

## Structure générale

PreLoRUGo utilise le *framework* iris écrit en go principalement pour le routage et le codage et décodage en json.

La structure s'inspire des principes du MVC avec une séparation des modèles et du traitement des requêtes. 

Un package configuration permet également de gérer des paramètres différents suivant qu'on est en test, en local ou sur le serveur de déploiement.

Le logiciel a été conçu pour être déployé sur un stack elastic beanstalk d'AWS avec un serveur go et une base de données PostgreSQL.

Le point d'entrée du serveur est le fichier main.go qui intègre les autres packages afin de récupérer la configuration, lancer la connexion à la base de données, lancer les migrations et lancer le serveur en renvoyant vers le package actions pour le traitement des requêtes.

PreLoRUGo comporte donc
* le fichier  `main.go` qui est le point d'entrée du fichier
* le package `config` qu contient toutes les fonctions de configuration pour lire `config.yml` contenant la configuration du serveur ou celle locale et des tests. Le package gère aussi la connexion à la base de données et la séquence d'initialisation. Celle-ci crée les tables de la base de données qui n'existent pas puis exécute les migrations pour modifier la structure des tables existantes
* le package `actions` qui contient le routage dans le fichier `routes.go` et l'ensemble des actions de traitement des requêtes qui appellent à leur tour les modèles pour récupérer les données depuis la base PostgreSQL et qui gère les erreurs
* le package `models` qui regroupe les modèles gérant les requêtes SQL pour récupérer les données depuis PostgreSQL

D'une manière générale les actions sont regroupées par fichier similaire à celui utilisé par le modèle. Par exemple, les requêtes de l'API relatives aux villes sont gérées par le fichier `city.go` du package `actions` et font appel au modèle `city` du package `models` qui comporte lui aussi un fichier `city.go`

## Fonctionnement des tests

Les tests sont implémentés uniquement pour les actions. Toutefois, la base de données ne fait l'objet d'aucun mocking. Les tests utilisent donc une base de données PostgreSQL de test dont les paramètres doivent être codés dans le fichier `config.yml`.

Les tests sont regroupés par une fonction de test qui s'assure de l'ordre de lancement des tests (certaines requêtes étant interdépendantes). Le fichier d'entrée des test est `commons_test.go` qui comporte le point d'entrée de tous les tests et qui assure la connexion spécifique à la base de données de test.

Cette base de test est configurée pour commencer par la suppression de toutes les tables et des views puis en appelant la fonction d'initialisation qui crée toutes les tables et lance les migrations. 

Cependant, ce processus ne permet de forcément de refléter le fonctionnement des migrations et la conception des tests doit en tenir compte notamment lorsque des colonnes sont ajoutées ou retirées de la base de données. L'ajout d'une colonne doit donc se faire à la fois dans la fonction de création des tables et dans la fonction qui gère les migrations, sinon les tests échoueront.
