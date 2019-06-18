#Kommunikation im ACPLT Server
Die urspr�ngliche Kommunikation des Systems l�uft �ber das ACPLT/KS Protokoll. Mittlerweile ist �ber die Bibliothek opcua aber auch ein OPCUA Server problemlos m�glich.

#Anwendungsbereich von ACPLT/KS
ACPLT/KS ist ein Kommunikationskonzept zum einfachen Austausch von aktuellen Zustandswerten, Modellinformationen und zum �bertragen von Modelltransformationsbefehlen. ACPLT/KS realisiert den Austausch in einzelnen Kommunikationsbeziehungen nach dem Client-Serverprinzip. ACPLT/KS betrachtet konzeptionell nur die einzelne Beziehung, l�sst aber beliebig viele voneinander unabh�ngige Beziehungen zu. In der Umsetzung kann ein Peer daher gleichzeitig mehrere Beziehungen unterhalten in denen er auch in unterschiedlichen Rollen (Client/Server) auftreten kann. 

In der einzelnen Beziehung realisiert ACPLT/KS eine 1:1 �Beziehung zwischen einem Client und einem Server. Bei jedem Datenaustausch schickt der Client einen Request und der Server antwortet mit genau einem Response, dann ist der Datenaustausch beendet. Der Client kann den n�chsten Request erst senden wenn der Response eingetroffen ist (au�er im Fehlerfall nach Timeout). Die Dienste die durch die ACPLT/KS Requests im Server angesto�en werden k�nnen als atomar betrachtet werden, d. h: ihre Ausf�hrung und die Generierung der Response erfolgt sofort und ohne Beeintr�chtigung der Echtzeit in der Serveranwendung. 

(Bemerkung: Dies ist durch die einfache Art der Dienste in fast allen �vern�nftigen� Anwendungsf�llen gegeben, zur Sicherstellung m�ssen allerdings Restriktionen wie Pfadtiefe, Anzahl der in einem Request gelesenen Variablen, Namensl�ngen usw. ber�cksichtigt werden. In der NAMUR-Empfehlung NE139 sind einige dieser Beschr�nkungen beschrieben. Es ist jedoch zu bemerken, dass in heutigen Systemen im Allgemeinen die Beschr�nkungen nicht zum Tragen kommen.)

Mit diesem Konzept ist die Kommunikation von Seiten des Klienten synchron und von Seiten des Servers zustandslos.

ACPLT/KS ist f�r viele Anwendungen ausreichend. Es erm�glicht insbesondere die Ankopplung an die Leitsysteme und den effizienten Datenaustausch zwischen Automatisierungskomponenten. F�r h�herwertige und komplexere Kommunikationsanforderungen dienen das ACPLT- Message-System (ACPLT/KS-M) und das ACPLT Service-System (ACPLT/KS-S). Beide k�nnen direkt auf das Basissystem ACPLT/KS aufgesetzt werden.

#Aufbau von ACPLT/KS

ACPLT/KS besteht aus drei Protokollschichten. Auf der obersten Schicht befindet sich die funktionale Semantik. Diese ist einheitlich festgelegt und charakterisiert die Kommunikationsfunktionalit�t von ACPLT/KS. Auf der darunterliegenden Schicht erfolgt die Kodierung der zu �bertragenden Informationen in eine Protokollsprache. Hier gibt es mehrere Varianten. Zurzeit werden ein bin�res Format (ksXDR) und eine textbasierte Kodierungen (ksHTTP) mit mehreren Varianten (plain / tcl / xml) unterst�tzt. Auf der untersten Ebene werden die ACPLT/KS Aufrufe auf die entsprechenden Bussysteme abgebildet. Zurzeit gibt es nur eine Abbildung auf TCP. Abbildungen z.B. auf CAN-Tunnel oder Profibus-Tunnel w�ren jedoch ebenfalls einfach realisierbar.

F�r viele Systeme (Leitsysteme, PIMS, Archivsysteme...) sind entsprechende Serverkomponenten vorhanden. (Die ACPLT/KS-Server sind konfigurationsfrei, sie m�ssen nur einmal installiert sein, dann steht die gesamte Funktionalit�t � auch bei nachtr�glichen �nderungen der Modelle im Zielsystem - immer vollst�ndig zur Verf�gung)


#Pfadsyntax in KS

Volle Pfade werden im KS folgenderma�en formuliert:

`//dev/ov_server_a/vendor/server_description`

Dies bezeichnet eine Variable namens /vendor/server_description auf einem Server namens ov_server_a dessen Zugriffsport man �ber einen KS-Manager auf Port 7509 unter dem Rechner-namen dev findet.

Das Adressziel ist das technologische Ziel und unabh�ngig von der Art der Kodierung (ksXDR und ksHTTP). Die Client-Anwendung kann w�hlen, �ber welches Protokoll sie den Server erreichen m�chte. Die Wege eines Zugriffs sind jedoch intern je nach gew�hltem Protokoll unterschiedlich.

Ist der Port des Servers bekannt, so kann die Abfrage beim KSMANAGER gespart und der Server direkt adressiert werden. Die Syntax lautet dann zum Beispiel f�r den Port 240678:

`//dev/ov_server_a:240678/vendor/server_description`

Eine Angabe des Ports des MANAGERS

`//dev:7509/ov_server_a:240678/vendor/server_description`

kann erfolgen, ist dann jedoch (wie der Servername selbst) ohne Bedeutung. Der Manager l�uft immer auf Port 7509!

G�ltige Zugriffe auf den Manager sind:

````
//dev:7509/MANAGER:7509/vendor/server_description
//dev/MANAGER:7509/vendor/server_description
//dev:7509/MANAGER/vendor/server_description
//dev/MANAGER/vendor/server_description
````

##ksXDR
Eine Anzeige aller vorhandenen Server wird �ber ein GetPP (nicht in kshttp realisiert) auf die Domain /servers des Managers (Port 7509 des Rechners) ausgel�st. Dabei wird davon ausgegangen, dass der Manager l�uft. Soll ein spezifischer Server angefragt werden, so kontaktiert die Applikation den Manager und setzt dort einen Befehl GETSERVER ab um den TCP-Port zu erhalten.

Ein Zugriff auf einen Server ohne Manager ist m�glich, aber wird von vielen Applikationen nicht unterst�tzt.

##ksHTTP
Eine Anzeige aller vorhandenen Server wird �ber ein GetEP auf die Domain /servers des Managers (Port 7509 des Rechners) ausgel�st. Dabei wird davon ausgegangen, dass der MANAGER l�uft. 

Soll ein spezifischer Server angefragt werden, so kontaktiert die Applikation den Manager und setzt dort einen Befehl getServer ab um den TCP-Port zu erhalten.

Wird ein Port in der Adresse angegeben, so wird kein MANAGER gefragt, sondern direkt eine Verbindung auf diesem Port aufgebaut.

###Beispiele f�r ksHTTP
`//dev/fb_control/vendor/server_description`

Da KEIN Port angegeben wurde muss der MANAGER (auf Port 7509) nach dem Port f�r fb_control gefragt werden

`//dev/fb_control:8080/vendor/server_description`

Da ein Port f�r den Server angegeben ist, greift das Programm direkt auf den Port 8080 zu und fragt keinen MANAGER um diesen zu erfahren.

Intern muss ein http-KS Client folgende http Ressourcen aufrufen:

````
http://dev:7509/getServers?servername=fb_control => 8080
http://dev:8080/getVar?path=vendor/server_description
````
