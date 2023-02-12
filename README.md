# **README**

**This program aims to work with goroutines which are functions running in the background concurrently.**

This program covers a real-world scenario, where a user wants to subscribe and buy one of a series of available subscriptions. 

When they purhcase a subscription, we will:
- generate an invoice, 
- send an email, 
- and generate a PDF manual,
- and send that to the user. 

We will do these things concurrently.

---

## **HOW TO START THE PROGRAM**

```bash
$docker compose up -d  # start docker containers in the background
$docker ps  # check running docker containers
$make start  # start the program using makefile 
```

## **HOW TO STOP THE PROGRAM**

```bash
$make stop  # stop the program using makefile
$docker compose down  # stop docker containers running in the background
$docker ps  # check running docker containers
```

---

## **USED LIBRARIES**
