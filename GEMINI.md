# Contexte Gemini : Projet Basics

- Ce document fournit le contexte pour le projet "Basics", un interpréteur BASIC pour vieux ordinateurs.
- Le projet se concentre actuellement principalement sur les ordinateurs APPLE II, avec une architecture permettant son extension à d'autres ordinateurs « rétro », par exemple Commodore, Oric-1, Zx Spectrum, Amstrad CPC-6128.

## Aperçu du projet

- **Nom** : Basics
- **Concept** : Interpréteur multi-BASIC et multi-architecture pour vieux ordinateurs
- **Genre** : Interpreteur, Langage BASIC
- **Plateforme** : Windows 10+, Lunix

## Stack Technologique

- **Langage** : Go 1.25+
- **Moteur de rendu graphique** : Ebiten 2.9.7
- **Outil de build** : Mage

## Structure du projet

```
/ basics
├── .vscode/                     # Paramètres spécifiques à l'éditeur VS Code
├── .asserts/                    # Fichiers d'assertions pour les tests de validation
├── .bintools/                   # Outils de gestion des polices binaires (ex: Excel pour Bitmap Font)
├── bin/                         # Dossier de destination des exécutables compilés
├── cmd/                         # Points d'entrée de l'application
│   └── basics/                  # Application principale : initialisation et boucle de l'interpréteur
├── examples/                    # Suite de tests et programmes BASIC d'exemple
│   ├── display/                 # Exemples liés à l'affichage (PRINT, HOME, FLASH...)
│   ├── flow_control/            # Exemples de structures de contrôle (FOR, IF, GOSUB...)
│   ├── input/                   # Exemples de gestion des entrées utilisateur (INPUT, GET)
│   ├── maths/                   # Exemples de fonctions mathématiques (SQR, INT, ABS...)
│   ├── operators/               # Exemples d'opérations booléennes et logiques
│   ├── programs/                # Sous-catégories de programmes complets
│   ├── release/                 # Algorithmes classiques (Primes, Fibonacci, Factorial...)
│   ├── strings/                 # Exemples de manipulation de chaînes (MID$, LEFT$...)
│   ├── tabs/                    # Exemples de positionnement (HTAB, VTAB)
│   └── variables/               # Exemples d'affectation et de typage
├── internal/                    # Cœur logique du projet (non exportable)
│   ├── app/                     # Orchestration du cycle de vie de l'application
│   ├── binary/                  # Manipulation de données et fichiers binaires
│   ├── common/                  # Utilitaires transverses partagés
│   ├── constants/               # Définitions des constantes globales du système
│   ├── errors/                  # Système de gestion et de formatage des erreurs BASIC
│   ├── input/                   # Gestion bas-niveau des entrées clavier/périphériques
│   ├── interpreter/             # Moteur d'exécution (interprétation des commandes)
│   ├── lexer/                   # Analyseur lexical (transformation du code source en tokens)
│   ├── logger/                  # Système de journalisation interne
│   ├── machines/                # Émulation des spécificités matérielles (Apple II, etc.)
|   |   ├── apple2/              # Implémentation complète de l'architecture Apple II (modes texte, palettes, mémoire)
|   |   └── tty/                 # Interface de type terminal standard (mode texte simple sans émulation graphique)
│   ├── parser/                  # Analyseur syntaxique (construction de l'arbre d'exécution)
│   ├── runtime/                 # Gestionnaire d'état (mémoire, variables, pile d'appels)
│   ├── token/                   # Définition des lexèmes et mots-clés du langage
│   └── video/                   # Moteur de rendu graphique et texte (via Ebiten)
|       ├── ebiten/              # Intégration et implémentation spécifique via la bibliothèque Ebiten
|       ├── font/                # Gestion des polices de caractères bitmap (7x8, 8x8) et générateurs de caractères
|       ├── renderer/            # Logique d'abstraction du rendu (interfaçage entre la logique vidéo et le moteur graphique)
|       └── text/                # Gestion des buffers de texte, des cellules de caractères et du défilement (scrolling)
├── release/                     # Artefacts et ressources pour la distribution
├── testutils/                   # Bibliothèques d'aide pour les tests unitaires et d'intégration
├── winres/                      # Ressources spécifiques à Windows (icônes de l'exécutable)
├── architecture.md              # Spécifications détaillées de l'architecture du projet
├── GEMINI.md                    # Instructions de contexte et règles pour l'assistant IA
├── magefile.go                  # Scripts d'automatisation du build (remplace Make)
└── go.mod                       # Fichier de définition du module Go et dépendances
```

## Conventions de codage

- Expert développeur Golang v1.25+
- Respect des bonnes pratiques de développement pour ce langage
- Le formatage du code doit respecter le standard issu de l'outil officiel `gofmt`.
- Utiliser le CamelCase pour le nommage des variables, structures, interfaces...
    - Une variable / fonction commençant par une minuscule -> privée
    - Une variable / fonction commençant par une majuscule -> publique
- Les noms des packages doivent être en minuscules, sans underscore et au singulier.
- Les noms des fichiers doivent être en minuscules avec possibilité de snake_case dans le nommage. Le nom des fichiers doit refléter la responsabilité qu'ils portent.
- Eviter tout Getter / Setter inutile.
- Toujour simplémenter une gestion explicite des erreurs.
- Tout nouveau code doit être documenté en langue anglaise.
- Tout package doit être documenté.
- Les commentaires multilignes sont autorisés.
- La documentation doit pouvoir être généré directement avec `godoc`.
- Pour les tests unitaires, utiliser des tests `table-driven`.
