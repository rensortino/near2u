select Dispositivo.name, Dispositivo.type, Dispositivo.code
 from ((Dispositivo join Sensore on Sensore.code = Dispositivo.code)
  join Dispositivo_Ambiente on Dispositivo.code = Dispositivo_Ambiente.code) 
  where Dispositivo_Ambiente.cod_ambiente = '"+ ambiente.getcodAmbiente() +"';




  select Dispositivo.name, Dispositivo.type, Dispositivo.code 
  from ((Dispositivo join Attuatore on Attuatore.code = Dispositivo.code) 
  join Dispositivo_Ambiente on Dispositivo.code = Dispositivo_Ambiente.code) 
  where Dispositivo_Ambiente.cod_ambiente = '"+ ambiente.getcodAmbiente() +"';


  select comando from Comandi where cod_attuatore = "+ code +";