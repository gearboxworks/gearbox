#Die textbasierte Kodierung ksHTTP 

Mit dem Aufkommen der SOA-Architekturen und des WSDL wurden alternativ zu den Bin�rformaten textbasierte Formate und Dienstprotokolle nach dem SOA-Konzept des W3C entwickelt. Diese Ans�tze haben sich jedoch nicht bew�hrt (extrem schlechte Performance) und wurden wieder eingestellt. Alternativ entstanden in den letzten Jahren mit den Web-Technologien schlanke textbasierte Protokolle die direkt in den Browsern generiert und in Servern mit integrierten Web-Diensten effizient verarbeitet werden k�nnen. ksHTTP ist ein Protokoll das den ACPLT-Standard auf einfache Weise auf ein textbasiertes Format umsetzt. ksHTTP unterst�tzt bis auf die Archivfunktionen alle Dienste der ACPLT/KS-Spezifikation.

ACPLT/KS unterst�tzt eine textbasierte Kodierung in drei Varianten:
1.	text/plain: Unterst�tzt eine einfache, jedoch strukturlose Kodierung
2.	text/xml unterst�tzt eine Kodierung in einem einfachen xml-format (KSX2.0-Spezifikation). 
3.	text/tcl unterst�tzt eine f�r tcl-Anwendungen optimierte Formatierung.
Die folgende Spezifikation beschreibt die Syntax der Aufrufe des Webservers.

#	URL Syntax

Alle Befehle von ACPLT/KS erlauben eine Parametrierung innerhalb des Requests. So k�nnen z. B. in einer Anfrage beliebig viele Parameter �bergeben werden. Zur Formulierung dieser Abh�ngigkeiten wird eine einfache Syntax in Anlehnung an den RFC 3986 benutzt. (http://tools.ietf.org/html/rfc3986#section-3.4). 

Nach dieser Spezifikation wird der Name vom Wert durch ein Gleichheitszeichen (=) getrennt, mehrere Name/Wert Paare werden jeweils durch ein Und-Zeichen (&) getrennt. Dies lehnt sich an die HTML Standards an, welche Formulare in �hnlicher Form (MIME Type application/x-www-form-urlencoded) �bermitteln. 

Alle Aufrufe werden im Server �ber Zugiffsfunktionen ausgef�hrt die die Erlaubnis und Korrektheit des Zugriffs pr�fen. Dazu werden im Header der Requests entsprechende �ticket� Informationen mit �bertragen. Dies sind hier nicht dargestellt.

##getServer
 
BasisURL: `/getServer`

Aufruf zum Abfragen des Servernamens

`/getServer?servername=fb_server_1`

##	register (wird nur intern genutzt)

`/register?name=fb_server_1&port=7510&version=2`

## unregister (wird nur intern genutzt)
`/unregister?name=fb_server_1&port=7510&version=2`

##	GetVar

BasisURL: `/getVar`

Aufruf zum Auslesen des Werts von Objektvariablen. Dies k�nnen sowohl Aktualwerte, Parameter, Eing�nge, Ausg�nge usw. sein. An den Objekten des Klassenmodells sind dies aber z.B. auch Typinformationen.

URL Parameter |Wert (Beispiel)|Beschreibung
------------- | -------------| -------------
path|/vendor/server_time|voller Pfad der Variable
stream|1|Erlaubt Server-Sent Events, welche im Standard http://dev.w3.org/html5/eventsource/ definiert sind.

Beispiele:
````
/getVar?path=/vendor/server_time
/getVar?path=/vendor/database_free&stream=1
````
GetVar erlaubt das Lesen mehrerer Variablen in einem Aufruf. Dazu wird jeder Pfad mit einer Nummer markiert:
`/getVar?path[0]=/TechUnits/add1.OUT&path[1]=/TechUnits/add2.OUT&path[2]= /TechUnits/add2.IN1`
Werden mehrere Parameter des gleichen Objekts direkt hintereinander gelesen, dann gen�gt es, den Objektpfad einmal beim ersten Parameter anzugeben. Die eben formulierte Abfrage kann also verk�rzt auch so formuliert werden:

Beispiel:

`/getVar?path[0]=/TechUnits/add1.OUT&path[1]=/TechUnits/add2.OUT&path[2]=.IN1`

##	SetVar
BasisURL: `/setVar`

Aufruf zum Setzen von Variablen. Mit dem Aufruf wird die Variable einmalig auf den mit�bertragenen Wert gesetzt.

URL Parameter |	Wert (Beispiel)	| Beschreibung
------------- | -------------| -------------
path|	/TechUnits/add1.IN1|	voller Pfad der Variable
newvalue|	1.5	|Neuer Wert der Variable
vartype	|KS_VT_BOOL	|Variablentyp wird mit gesetzt (optional bei any-Eing�ngen)

Beispiel:

`/setVar?path=/TechUnits/add1.IN1&newvalue=4.5` 

SetVar erlaubt das Setzen mehrerer Variablen in einem Aufruf. Dazu wird jeder Pfad mit einer Nummer markiert:

`/setVar?path[0]=/TechUnits/add1.IN1&path[1]=/TechUnits/add2.IN1&path[2]=/TechUnits/add2.IN2&newvalue[0]=3.2&newvalue[1]=2&newvalue[2]=1`

Werden mehrere Parameter des gleichen Objekts direkt hintereinander gesetzt, dann gen�gt es, den Objektpfad einmal beim ersten Parameter anzugeben.

Beispiel:

`/setVar?path[0]=/TechUnits/add1.IN1&path[1]=/TechUnits/add2.IN1&path[2]=.IN2&newvalue[0]=3.2&newvalue[1]=2&newvalue[2]=1`

Hinweis: Die Werte m�ssen URLencodet �bertragen werden (zum Beispiel �ber die JavaScriptfunktion encodeURIComponent). Einzelne Eintr�ge von Vektoren werden in geschweifte Klammern gekapselt und per Leerzeichen getrennt: �{hello} {world}�

##getEP
BasisURL: `/getEP`

URL Parameter| 	Wert (Beispiel)	|Beschreibung
------------- | -------------| -------------
path|	/TechUnits/add1	|voller Pfad des Objektes dessen Kinder gepr�ft werden
requestType	|OT_VARIABLE	|Schr�nkt den Typ der Ergebnisse ein (optional)
requestOutput	|OP_NAME	|
nameMask	|*|	erlaubt eine Suche nach Objekten
scopeFlags	|�parts� oder �children�	|Nur Parts oder nur Kinder werden abgefragt

Beispiel:

`/getEP?path=/TechUnits/add1`

requestType kann OT_DOMAIN (default), OT_VARIABLE, OT_LINK oder OT_ANY sein. OT_DOMAIN liefert eine Liste der Kinder der Assoziation ov_containment, OT_VARIABLE alle Variablen eines Objekts und OT_LINK eine Liste aller Assoziationen. OT_ANY liefert alle drei Typen.

requestOutput (bzw mehrere requestOutput[?]) liefert eine beliebige Anzahl von Informationen zum Ergebnis. Wird OP_ANY angefragt, so ist die Reihenfolge wie in Spalte 2 beschrieben:

Schl�sselwort|	OP_ANY	|Beschreibung
------------- | -------------| -------------
OP_NAME|	0|	Name des Objekts
OP_TYPE|	1|	Typ der Antwort (KS_OT_DOMAIN, KS_OT_VARIABLE oder KS_OT_LINK)
OP_COMMENT	|2|	Kommentar des Objekts
OP_ACCESS	|3	|Eine per Leerzeichen getrennte Liste von Zugriffsrechten, die der Benutzer hat: KS_AC_NONE, KS_AC_READ, KS_AC_WRITE, KS_AC_DELETEABLE, KS_AC_INSTANTIABLE, KS_AC_UNLINKABLE, KS_AC_LINKABLE, KS_AC_RENAMEABLE
OP_SEMANTIC	|4|	ASCII Repr�sentanz aller gesetzten Flags des Objektes. Alphabe-tisch, ohne Trenner
OP_CREATIONTIME	|5|	Erstellungsdatum des Objekts
OP_CLASS	|6|	Im Falle von KS_OT_DOMAIN Klasse des Objekts, bei KS_OT_VARIABLE Typ der Variable wie KS_VT_BOOL, in KS_OT_LINK eine Antwort wie (KS_LT_LOCAL_1_1, KS_LT_LOCAL_1_MANY, KS_LT_LOCAL_MANY_MANY, KS_LT_LOCAL_MANY_1, KS_LT_GLOBAL_1_1...)
OP_ASSOCIDENT	||	Im Falle eines Links die zugeh�rige Assoziation
OP_ROLEIDENT	||	Im Falle eines Links die zugeh�rige Rolle
OP_TECHUNIT	||	Im Falle einer Variable: physikalische Einheit der Variable als String

Hinweis: Da die Kodierung nach KSX eine einfache Interpretation der Werte erm�glicht werden dort immer alle Informationen (�quivalent zu OP_ANY) geliefert.

nameMask erlaubt eine Suche nach Objektnamen. Default ist �*�, also eine Suche nach allen Namen. �*� gilt f�r beliebige Anzahl von beliebigen Zeichen, �?� steht f�r genau ein Zeichen.

##createObject
BasisURL: `/createObject`

Request zum dynamischen Erstellen eines Objekts im Zielsystem. Zugrundeliegende Vorstellung: Zu jedem Typ gibt es ein Objekt das als Factory die Erstellung der Instanzen des Typs realisiert.

URL Parameter |	Wert (Beispiel)|	Beschreibung
------------- | -------------| -------------
path|	/TechUnits/add1|	voller Pfad des neuen Objekts
factory	|/acplt/iec61131stdfb/add|	Klassenobjekt des neuen Objekts

Beispiel:
`/createObject?path=/TechUnits/add1&factory=/acplt/iec61131stdfb/add`

createObject erlaubt das Anlagen mehrerer Objekte auf einmal. Werden mehrere Objekte des gleichen Typs direkt hintereinander erstellt, dann gen�gt es die Factory-Information an den letzten Eintrag dieser Reihe anzuh�ngen. 

`/createObject?Path[1]=/TechUnits/add1&path[2]=/TechUnits/add2&factory=/acplt/iec61131stdfb/add`

##	deleteObject
BasisURL: `/deleteObject`

Request zum dynamischen L�schen eines Objekts im Zielsystem.

URL Parameter |	Wert (Beispiel)	|Beschreibung
------------- | -------------| -------------
path	|/TechUnits/add1|	voller Pfad des zu l�schenden Objekts

Beispiel:
`/deleteObject?path=/TechUnits/add1`

deleteObject erlaubt das L�schen mehrerer Objekte auf einmal. Dabei wird der Parameter entsprechend angepasst. 

`/deleteObject?path[1]=/TechUnits/add1&path[2]=/TechUnits/add2`

##	renameObject

BasisURL: `/renameObject`
Request zum dynamischen �ndern des Namens eines Objekts im Zielsystem.

URL Parameter |	Wert (Beispiel)	|Beschreibung
------------- | -------------| -------------
path|	/TechUnits/add1	|voller Pfad des Objekts
newname|	/TechUnits/add2|	voller neuer Pfad des Objekts

Beispiel:

`/renameObject?path=/TechUnits/add1&newname=/TechUnits/add2`

renameObject erlaubt das Umbenennen mehrerer Objekte auf einmal. Dabei wird der Parameter entsprechend angepasst. 
`/renameObject?path[1]=/TechUnits/add1&path[2]=/TechUnits/sub1&newname[1]=/TechUnits/add2&newname[2]=/TechUnits/sub2`

##	link
BasisURL: `/link`

Request zum dynamischen Verlinken eines Objekts im Zielsystem.

URL Parameter |	Wert (Beispiel)	|Beschreibung
------------- | -------------| -------------
path	|/TechUnits/add1.taskparent	|voller Pfad des Objekts inkl Assoziationsname
element|	/UrTask|	voller Pfad des Linkziels
placehint	|BEGIN|	Placementhint (DEFAULT=END, BEGIN, END, AFTER, BEFORE)
placepath	|/TechUnits/sub1|	wird benutzt als Referenz wenn BEFORE oder AFTER
oppositeplacehint	|AFTER|	Zweiter Placementhint bei m-n-Assoziationen
oppositeplacepath|	/TechUnits/abs1|	Zweiter placepath bei m-n-Assoziationen

Beispiel:
`/link?path=/TechUnits/add1.taskparent&element=/UrTask`

Link erlaubt das Linken mehrerer Objekte auf einmal. Dabei werden die Parameter entsprechend angepasst. 

##	unlink
BasisURL: `/unlink`

Request zum dynamischen L�schen eines Links.

URL Parameter |	Wert (Beispiel)	|Beschreibung
------------- | -------------| -------------
path	|/TechUnits/add1.taskparent|	voller Pfad des Objekts inkl Assoziationsname
element	|/UrTask|	voller Pfad des Linkziels

Beispiel:
`/unlink?path=/TechUnits/add1.taskparent&element=/UrTask`

Unlink erlaubt das Unlinken mehrerer Objekte auf einmal. Dabei wird der Parameter entsprechend angepasst. 
