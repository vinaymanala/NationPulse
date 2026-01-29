# confirm topics exist (on the host)

docker exec -it kafka-1 /opt/kafka/bin/kafka-topics.sh --list --bootstrap-server kafka-1:9092

# watch bff/reporting logs for errors

docker compose logs -f nationpulse-bff nationpulse-reporting-svc

# pg functions

CREATE OR REPLACE FUNCTION public.geteconomygdpbycountry(p*country_code character varying)
RETURNS SETOF egdptable
LANGUAGE plpgsql
AS $function$
/* gdppercapitads \_/
BEGIN
RETURN QUERY
SELECT \* FROM egdptable
WHERE country_code LIKE p_country_code;  
END;
$function$

CREATE OR REPLACE FUNCTION public.geteconomygovenmentbycountry(p*country_code character varying)
RETURNS SETOF egovtable
LANGUAGE plpgsql
AS $function$
/* publicgovernmentyearlyds \_/
BEGIN
RETURN QUERY
SELECT \* FROM egovtable
WHERE country_code LIKE p_country_code;
END;
$function$

CREATE OR REPLACE FUNCTION public.getgrowthgdpofcountry(p*country_code character varying)
RETURNS SETOF ggdptable
LANGUAGE plpgsql
AS $function$
/* perfgrowthgdpds \_/
BEGIN
RETURN QUERY
SELECT \* FROM ggdptable WHERE country_code LIKE p_country_code;

END;
$function$

CREATE OR REPLACE FUNCTION public.getgrowthpopnbycountry(p*country_code character varying)
RETURNS SETOF gpopulationtable
LANGUAGE plpgsql
AS $function$
/* perfgrowthpopulationds \_/
BEGIN
RETURN QUERY
SELECT \* FROM gpopulationtable WHERE country_code LIKE p_country_code;

END;
$function$

CREATE OR REPLACE FUNCTION public.gethealthcasesbycountry(p*country_code character varying)
RETURNS SETOF htable
LANGUAGE plpgsql
AS $function$
/* healthstatusds \_/
BEGIN
RETURN QUERY
SELECT \* FROM htable
WHERE country_code LIKE p_country_code;
END;
$function$

CREATE OR REPLACE FUNCTION public.getpopulationofcountry(p*country_code character varying)
RETURNS SETOF poptable
LANGUAGE plpgsql
AS $function$
/* populationds \_/
BEGIN
RETURN QUERY
SELECT \* FROM poptable
WHERE country_code LIKE p_country_code;
END;
$function$
