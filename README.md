# elrek-system_GO

Ez lenni repository for elrek-system_GO

Ez a projekt lenne a projektmunka a Felhő és Devops gyakorlatra.

## Ami meg lett valósítva, mint követelmény ebben a konkrét projektben:

### Jenkins
A projekthez tartoznak tesztek, jelen esetben az endpointokat és azok műveleteiket összességében tesztelő E2E tesztsorról van szó GO nyelven.

Egy saját szerveren futó Jenkins instance-ben multibranch jobként lett felvéve a GitHub repositoryja a projektnek, további beállítások után már érzékeli is a meglévő brancheket és pull requesteket. Emiatt úgy lett beállítva a Jenkinsfile, hogy minden pull request nyitáskor lefuttatja a teszteket, majd csak a dev és main branchre pusholás esetén fut le a projekt és a Docker image buildelése is. 

### Terraform
Miután felkerült a Dockerhubra felkerült a legfrissebb image, majd törölte is a szerveren futó meglévő konténereket és image-eket, egy Terraform fájl segítségével újra deployolásra kerül a legfrissebb verzió a megfelelő környezeti változókkal

### Graylog
Ezek mellett az API rá lett kötve egy másik szerveren futó Graylog instance-ra logoláshoz. 

### Nginx Proxy Manager
Innentől nem a konkrét projektben lettek megvalósítva a következő toolok, hanem a teljes rendszer működéséhez szükségesek. Az egyik egy reverse proxy az API biztonságos eléréséhez, ehhez a következő open-source projekt fut konténerként:

### Grafana
További teljes rendszerhez szükséges rendszer a monitorozáshoz használt Grafana, annyi különbséggel, hogy nem Prometheus van mögötte, hanem InfluxDB, illetve nem az API-t figyeli, hanem az őt futtató szervereket monitorozza.

### Üzembe helyezés
Ezt az összes rendszert jelenleg nem lehet 1 terraform futtatásával elindítnai és üzembe helyezni, ugyanis jelen esetben ezek különböző szervereken futnak és különböző rendszerekhez is hozzá vannak közve