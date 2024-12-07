# elrek-system_GO

Ez lenni repository for elrek-system_GO

Ez a projekt lenne a projektmunka a Felhő és Devops gyakorlatra.

## Ami meg lett valósítva, mint követelmény ebben a konkrét projektben:

### Jenkins
A projekthez tartoznak tesztek, jelen esetben az endpointokat és azok műveleteiket összességében tesztelő E2E tesztsorról van szó GO nyelven.

Egy saját szerveren futó Jenkins instance-ben multibranch jobként lett felvéve a GitHub repositoryja a projektnek, további beállítások után már érzékeli is a meglévő brancheket és pull requesteket. Emiatt úgy lett beállítva a Jenkinsfile, hogy minden pull request nyitáskor lefuttatja a teszteket, majd csak a dev és main branchre pusholás esetén fut le a projekt és a Docker image buildelése is. <img width="1350" alt="Screenshot 2024-12-07 at 2 40 40 PM" src="https://github.com/user-attachments/assets/3a3ef014-4a6a-45c1-964d-33526048f186">


### Terraform
Miután felkerült a Dockerhubra felkerült a legfrissebb image, majd törölte is a szerveren futó meglévő konténereket és image-eket, egy Terraform fájl segítségével újra deployolásra kerül a legfrissebb verzió a megfelelő környezeti változókkal

### Graylog
Ezek mellett az API rá lett kötve egy másik szerveren futó Graylog instance-ra logoláshoz. <img width="1440" alt="Screenshot 2024-12-06 at 11 07 33 PM" src="https://github.com/user-attachments/assets/99bcda0a-7245-414e-b3e6-94f5e5693e6f">


### Nginx Proxy Manager
Innentől nem a konkrét projektben lettek megvalósítva a következő toolok, hanem a teljes rendszer működéséhez szükségesek. Az egyik egy reverse proxy az API biztonságos eléréséhez, ehhez a következő open-source projekt fut konténerként:<img width="831" alt="Screenshot 2024-12-07 at 2 48 01 PM" src="https://github.com/user-attachments/assets/173146e6-49f4-483e-81eb-782302393c3f">
<img width="1440" alt="Screenshot 2024-12-06 at 11 07 22 PM" src="https://github.com/user-attachments/assets/0aae3124-f543-4268-95cc-687518f69373">


### Grafana
További teljes rendszerhez szükséges rendszer a monitorozáshoz használt Grafana, annyi különbséggel, hogy nem Prometheus van mögötte, hanem InfluxDB, illetve nem az API-t figyeli, hanem az őt futtató szervereket monitorozza.<img width="1440" alt="Screenshot 2024-12-06 at 11 11 28 PM" src="https://github.com/user-attachments/assets/2c3abe3f-1d7d-426a-90ca-5a71ef857363">
<img width="1428" alt="Screenshot 2024-12-06 at 11 03 21 PM" src="https://github.com/user-attachments/assets/699cfb46-f85e-4703-aa87-82032f2f8d48">


### Üzembe helyezés
Ezt az összes rendszert jelenleg nem lehet 1 terraform futtatásával elindítnai és üzembe helyezni, ugyanis jelen esetben ezek különböző szervereken futnak és különböző rendszerekhez is hozzá vannak közve
