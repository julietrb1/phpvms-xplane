package models

import "time"

type fetch struct {
	Userid   string        `json:"userid"`
	StaticId FlexibleField `json:"static_id"`
	Status   string        `json:"status"`
	Time     string        `json:"time"`
}

type params struct {
	RequestId     string        `json:"request_id"`
	SequenceId    string        `json:"sequence_id"`
	StaticId      FlexibleField `json:"static_id"`
	UserId        string        `json:"user_id"`
	TimeGenerated string        `json:"time_generated"`
	XmlFile       string        `json:"xml_file"`
	OfpLayout     string        `json:"ofp_layout"`
	Airac         string        `json:"airac"`
	Units         string        `json:"units"`
}

type general struct {
	Release           string        `json:"release"`
	ICAOAirline       FlexibleField `json:"icao_airline"`
	FlightNumber      string        `json:"flight_number"`
	IsEtops           string        `json:"is_etops"`
	DxRmk             string        `json:"dx_rmk"`
	SysRmk            FlexibleField `json:"sys_rmk"`
	IsDetailedProfile string        `json:"is_detailed_profile"`
	CruiseProfile     string        `json:"cruise_profile"`
	ClimbProfile      string        `json:"climb_profile"`
	DescentProfile    string        `json:"descent_profile"`
	AlternateProfile  string        `json:"alternate_profile"`
	ReserveProfile    string        `json:"reserve_profile"`
	Costindex         string        `json:"costindex"`
	ContRule          string        `json:"cont_rule"`
	InitialAltitude   string        `json:"initial_altitude"`
	StepclimbString   string        `json:"stepclimb_string"`
	AvgTempDev        string        `json:"avg_temp_dev"`
	AvgTropopause     string        `json:"avg_tropopause"`
	AvgWindComp       string        `json:"avg_wind_comp"`
	AvgWindDir        string        `json:"avg_wind_dir"`
	AvgWindSpd        string        `json:"avg_wind_spd"`
	GcDistance        string        `json:"gc_distance"`
	RouteDistance     string        `json:"route_distance"`
	AirDistance       string        `json:"air_distance"`
	TotalBurn         string        `json:"total_burn"`
	CruiseTAS         string        `json:"cruise_tas"`
	CruiseMach        string        `json:"cruise_mach"`
	Passengers        string        `json:"passengers"`
	Route             string        `json:"route"`
	RouteIfps         string        `json:"route_ifps"`
	RouteNavigraph    string        `json:"route_navigraph"`
	SidIdent          FlexibleField `json:"sid_ident"`
	SidTrans          FlexibleField `json:"sid_trans"`
	StarIdent         FlexibleField `json:"star_ident"`
	StarTrans         FlexibleField `json:"star_trans"`
}

type airport struct {
	ICAOCode        FlexibleField `json:"icao_code"`
	IATACode        FlexibleField `json:"iata_code"`
	FAACode         FlexibleField `json:"faa_code"`
	ICAORegion      string        `json:"icao_region"`
	Elevation       string        `json:"elevation"`
	PosLat          string        `json:"pos_lat"`
	PosLong         string        `json:"pos_long"`
	Name            string        `json:"name"`
	Timezone        string        `json:"timezone"`
	PlanRwy         string        `json:"plan_rwy"`
	TransAlt        string        `json:"trans_alt"`
	TransLevel      string        `json:"trans_level"`
	METAR           string        `json:"metar"`
	METARTime       FlexibleField `json:"metar_time"`
	METARCategory   FlexibleField `json:"metar_category"`
	METARVisibility FlexibleField `json:"metar_visibility"`
	METARCeiling    FlexibleField `json:"metar_ceiling"`
	TAF             string        `json:"taf"`
	TAFTime         FlexibleField `json:"taf_time"`
	ATIS            FlexibleField `json:"atis"`
}

type alternateAirport struct {
	ICAOCode        FlexibleField `json:"icao_code"`
	IATACode        FlexibleField `json:"iata_code"`
	FAACode         FlexibleField `json:"faa_code"`
	ICAORegion      string        `json:"icao_region"`
	Elevation       string        `json:"elevation"`
	PosLat          string        `json:"pos_lat"`
	PosLong         string        `json:"pos_long"`
	Name            string        `json:"name"`
	Timezone        string        `json:"timezone"`
	PlanRwy         string        `json:"plan_rwy"`
	TransAlt        string        `json:"trans_alt"`
	TransLevel      string        `json:"trans_level"`
	CruiseAltitude  string        `json:"cruise_altitude"`
	Distance        string        `json:"distance"`
	GcDistance      string        `json:"gc_distance"`
	AirDistance     string        `json:"air_distance"`
	TrackTrue       string        `json:"track_true"`
	TrackMag        string        `json:"track_mag"`
	TAS             string        `json:"tas"`
	Gs              string        `json:"gs"`
	AvgWindComp     string        `json:"avg_wind_comp"`
	AvgWindDir      string        `json:"avg_wind_dir"`
	AvgWindSpd      string        `json:"avg_wind_spd"`
	AvgTropopause   string        `json:"avg_tropopause"`
	AvgTdv          string        `json:"avg_tdv"`
	ETE             string        `json:"ete"`
	Burn            string        `json:"burn"`
	Route           string        `json:"route"`
	RouteIfps       string        `json:"route_ifps"`
	METAR           string        `json:"metar"`
	METARTime       time.Time     `json:"metar_time"`
	METARCategory   string        `json:"metar_category"`
	METARVisibility string        `json:"metar_visibility"`
	METARCeiling    string        `json:"metar_ceiling"`
	TAF             string        `json:"taf"`
	TAFTime         time.Time     `json:"taf_time"`
	ATIS            FlexibleField `json:"atis"`
}

type navlogFix struct {
	Ident           string        `json:"ident"`
	Name            string        `json:"name"`
	Type            string        `json:"type"`
	ICAORegion      interface{}   `json:"icao_region"`
	RegionCode      interface{}   `json:"region_code"`
	Frequency       FlexibleField `json:"frequency"`
	PosLat          string        `json:"pos_lat"`
	PosLong         string        `json:"pos_long"`
	Stage           string        `json:"stage"`
	ViaAirway       string        `json:"via_airway"`
	IsSidStar       string        `json:"is_sid_star"`
	Distance        string        `json:"distance"`
	TrackTrue       string        `json:"track_true"`
	TrackMag        string        `json:"track_mag"`
	HeadingTrue     string        `json:"heading_true"`
	HeadingMag      string        `json:"heading_mag"`
	AltitudeFeet    string        `json:"altitude_feet"`
	IndAirspeed     string        `json:"ind_airspeed"`
	TrueAirspeed    string        `json:"true_airspeed"`
	Mach            string        `json:"mach"`
	MachThousandths string        `json:"mach_thousandths"`
	WindComponent   string        `json:"wind_component"`
	Groundspeed     string        `json:"groundspeed"`
	TimeLeg         string        `json:"time_leg"`
	TimeTotal       string        `json:"time_total"`
	FuelFlow        string        `json:"fuel_flow"`
	FuelLeg         string        `json:"fuel_leg"`
	FuelTotalused   string        `json:"fuel_totalused"`
	FuelMinOnboard  string        `json:"fuel_min_onboard"`
	FuelPlanOnboard string        `json:"fuel_plan_onboard"`
	Oat             string        `json:"oat"`
	OatISADev       string        `json:"oat_isa_dev"`
	WindDir         string        `json:"wind_dir"`
	WindSpd         string        `json:"wind_spd"`
	Shear           string        `json:"shear"`
	TropopauseFeet  string        `json:"tropopause_feet"`
	GroundHeight    string        `json:"ground_height"`
	Fir             string        `json:"fir"`
	FirUnits        string        `json:"fir_units"`
	FirValidLevels  string        `json:"fir_valid_levels"`
	MORA            string        `json:"mora"`
	WindData        struct {
		Level []struct {
			Altitude string `json:"altitude"`
			WindDir  string `json:"wind_dir"`
			WindSpd  string `json:"wind_spd"`
			Oat      string `json:"oat"`
		} `json:"level"`
	} `json:"wind_data"`
	FirCrossing FlexibleField `json:"fir_crossing"`
}

type runway struct {
	Identifier         string        `json:"identifier"`
	Length             string        `json:"length"`
	LengthTora         string        `json:"length_tora"`
	LengthToda         string        `json:"length_toda"`
	LengthAsda         string        `json:"length_asda"`
	LengthLda          string        `json:"length_lda"`
	Elevation          string        `json:"elevation"`
	Gradient           string        `json:"gradient"`
	TrueCourse         string        `json:"true_course"`
	MagneticCourse     string        `json:"magnetic_course"`
	HeadwindComponent  string        `json:"headwind_component"`
	CrosswindComponent string        `json:"crosswind_component"`
	IlsFrequency       FlexibleField `json:"ils_frequency"`
	FlapSetting        string        `json:"flap_setting"`
	ThrustSetting      string        `json:"thrust_setting"`
	BleedSetting       string        `json:"bleed_setting"`
	AntiIceSetting     string        `json:"anti_ice_setting"`
	FlexTemperature    FlexibleField `json:"flex_temperature"`
	MaxTemperature     interface{}   `json:"max_temperature"`
	MaxWeight          interface{}   `json:"max_weight"`
	LimitCode          interface{}   `json:"limit_code"`
	LimitObstacle      FlexibleField `json:"limit_obstacle"`
	SpeedsV1           interface{}   `json:"speeds_v1"`
	SpeedsVr           interface{}   `json:"speeds_vr"`
	SpeedsV2           interface{}   `json:"speeds_v2"`
	SpeedsV2Id         string        `json:"speeds_v2_id"`
	SpeedsOther        FlexibleField `json:"speeds_other"`
	SpeedsOtherId      FlexibleField `json:"speeds_other_id"`
	DistanceDecide     string        `json:"distance_decide"`
	DistanceReject     string        `json:"distance_reject"`
	DistanceMargin     string        `json:"distance_margin"`
	DistanceContinue   string        `json:"distance_continue"`
}

type runwayDistance struct {
	Weight           string `json:"weight"`
	FlapSetting      string `json:"flap_setting"`
	BrakeSetting     string `json:"brake_setting"`
	ReverserCredit   string `json:"reverser_credit"`
	SpeedsVref       string `json:"speeds_vref"`
	ActualDistance   string `json:"actual_distance"`
	FactoredDistance string `json:"factored_distance"`
}

type landingRunway struct {
	Identifier         string        `json:"identifier"`
	Length             string        `json:"length"`
	LengthTora         string        `json:"length_tora"`
	LengthToda         string        `json:"length_toda"`
	LengthAsda         string        `json:"length_asda"`
	LengthLda          string        `json:"length_lda"`
	Elevation          string        `json:"elevation"`
	Gradient           string        `json:"gradient"`
	TrueCourse         string        `json:"true_course"`
	MagneticCourse     string        `json:"magnetic_course"`
	HeadwindComponent  string        `json:"headwind_component"`
	CrosswindComponent string        `json:"crosswind_component"`
	IlsFrequency       FlexibleField `json:"ils_frequency"`
	MaxWeightDry       string        `json:"max_weight_dry"`
	MaxWeightWet       FlexibleField `json:"max_weight_wet"`
}

type aircraft struct {
	ICAOcode         string        `json:"icaocode"`
	IATAcode         FlexibleField `json:"iatacode"`
	BaseType         string        `json:"base_type"`
	ICAOCode         string        `json:"icao_code"`
	IATACode         FlexibleField `json:"iata_code"`
	Name             string        `json:"name"`
	Reg              string        `json:"reg"`
	Fin              string        `json:"fin"`
	Selcal           string        `json:"selcal"`
	Equip            string        `json:"equip"`
	EquipCategory    string        `json:"equip_category"`
	EquipNavigation  string        `json:"equip_navigation"`
	EquipTransponder string        `json:"equip_transponder"`
	Fuelfact         string        `json:"fuelfact"`
	Fuelfactor       string        `json:"fuelfactor"`
	MaxPassengers    string        `json:"max_passengers"`
	SupportsTlr      string        `json:"supports_tlr"`
	InternalId       string        `json:"internal_id"`
	IsCustom         string        `json:"is_custom"`
}

type SimBriefOFP struct {
	Fetch       fetch            `json:"fetch"`
	Params      params           `json:"params"`
	General     general          `json:"general"`
	Origin      airport          `json:"origin"`
	Destination airport          `json:"destination"`
	Alternate   alternateAirport `json:"alternate"`
	TakeoffAltn FlexibleField    `json:"takeoff_altn"`
	EnrouteAltn FlexibleField    `json:"enroute_altn"`
	Navlog      struct {
		Fix []navlogFix `json:"fix"`
	} `json:"navlog"`
	Etops FlexibleField `json:"etops"`
	Tlr   struct {
		Takeoff struct {
			Conditions struct {
				AirportICAO      string `json:"airport_icao"`
				PlannedRunway    string `json:"planned_runway"`
				PlannedWeight    string `json:"planned_weight"`
				WindDirection    string `json:"wind_direction"`
				WindSpeed        string `json:"wind_speed"`
				Temperature      string `json:"temperature"`
				Altimeter        string `json:"altimeter"`
				SurfaceCondition string `json:"surface_condition"`
			} `json:"conditions"`
			Runway []runway `json:"runway"`
		} `json:"takeoff"`
		Landing struct {
			Conditions struct {
				AirportICAO      string `json:"airport_icao"`
				PlannedRunway    string `json:"planned_runway"`
				PlannedWeight    string `json:"planned_weight"`
				FlapSetting      string `json:"flap_setting"`
				WindDirection    string `json:"wind_direction"`
				WindSpeed        string `json:"wind_speed"`
				Temperature      string `json:"temperature"`
				Altimeter        string `json:"altimeter"`
				SurfaceCondition string `json:"surface_condition"`
			} `json:"conditions"`
			DistanceDry runwayDistance  `json:"distance_dry"`
			DistanceWet runwayDistance  `json:"distance_wet"`
			Runway      []landingRunway `json:"runway"`
		} `json:"landing"`
	} `json:"tlr"`
	Atc struct {
		FlightplanText string        `json:"flightplan_text"`
		Route          string        `json:"route"`
		RouteIfps      string        `json:"route_ifps"`
		Callsign       string        `json:"callsign"`
		FlightType     string        `json:"flight_type"`
		FlightRules    string        `json:"flight_rules"`
		InitialSpd     string        `json:"initial_spd"`
		InitialSpdUnit string        `json:"initial_spd_unit"`
		InitialAlt     string        `json:"initial_alt"`
		InitialAltUnit string        `json:"initial_alt_unit"`
		Section18      string        `json:"section18"`
		FirOrig        string        `json:"fir_orig"`
		FirDest        string        `json:"fir_dest"`
		FirAltn        string        `json:"fir_altn"`
		FirEtops       FlexibleField `json:"fir_etops"`
		FirEnroute     FlexibleField `json:"fir_enroute"`
	} `json:"atc"`
	Aircraft aircraft `json:"aircraft"`
	Fuel     struct {
		Taxi          string `json:"taxi"`
		EnrouteBurn   string `json:"enroute_burn"`
		Contingency   string `json:"contingency"`
		AlternateBurn string `json:"alternate_burn"`
		Reserve       string `json:"reserve"`
		Etops         string `json:"etops"`
		Extra         string `json:"extra"`
		ExtraRequired string `json:"extra_required"`
		ExtraOptional string `json:"extra_optional"`
		MinTakeoff    string `json:"min_takeoff"`
		PlanTakeoff   string `json:"plan_takeoff"`
		PlanRamp      string `json:"plan_ramp"`
		PlanLanding   string `json:"plan_landing"`
		AvgFuelFlow   string `json:"avg_fuel_flow"`
		MaxTanks      string `json:"max_tanks"`
	} `json:"fuel"`
	FuelExtra struct {
		Bucket []struct {
			Label    string      `json:"label"`
			Fuel     string      `json:"fuel"`
			Time     string      `json:"time"`
			Required interface{} `json:"required"`
		} `json:"bucket"`
	} `json:"fuel_extra"`
	Times struct {
		EstTimeEnroute   string `json:"est_time_enroute"`
		SchedTimeEnroute string `json:"sched_time_enroute"`
		SchedOut         string `json:"sched_out"`
		SchedOff         string `json:"sched_off"`
		SchedOn          string `json:"sched_on"`
		SchedIn          string `json:"sched_in"`
		SchedBlock       string `json:"sched_block"`
		EstOut           string `json:"est_out"`
		EstOff           string `json:"est_off"`
		EstOn            string `json:"est_on"`
		EstIn            string `json:"est_in"`
		EstBlock         string `json:"est_block"`
		OrigTimezone     string `json:"orig_timezone"`
		DestTimezone     string `json:"dest_timezone"`
		TaxiOut          string `json:"taxi_out"`
		TaxiIn           string `json:"taxi_in"`
		ReserveTime      string `json:"reserve_time"`
		Endurance        string `json:"endurance"`
		ContfuelTime     string `json:"contfuel_time"`
		EtopsfuelTime    string `json:"etopsfuel_time"`
		ExtrafuelTime    string `json:"extrafuel_time"`
	} `json:"times"`
	Weights struct {
		Oew            string `json:"oew"`
		PaxCount       string `json:"pax_count"`
		BagCount       string `json:"bag_count"`
		PaxCountActual string `json:"pax_count_actual"`
		BagCountActual string `json:"bag_count_actual"`
		PaxWeight      string `json:"pax_weight"`
		BagWeight      string `json:"bag_weight"`
		FreightAdded   string `json:"freight_added"`
		Cargo          string `json:"cargo"`
		Payload        string `json:"payload"`
		EstZfw         string `json:"est_zfw"`
		MaxZfw         string `json:"max_zfw"`
		EstTow         string `json:"est_tow"`
		MaxTow         string `json:"max_tow"`
		MaxTowStruct   string `json:"max_tow_struct"`
		TowLimitCode   string `json:"tow_limit_code"`
		EstLdw         string `json:"est_ldw"`
		MaxLdw         string `json:"max_ldw"`
		EstRamp        string `json:"est_ramp"`
	} `json:"weights"`
	Impacts struct {
		Minus6000Ft struct {
			TimeEnroute    string `json:"time_enroute"`
			TimeDifference string `json:"time_difference"`
			EnrouteBurn    string `json:"enroute_burn"`
			BurnDifference string `json:"burn_difference"`
			RampFuel       string `json:"ramp_fuel"`
			InitialFl      string `json:"initial_fl"`
			InitialTAS     string `json:"initial_tas"`
			InitialMach    string `json:"initial_mach"`
			CostIndex      string `json:"cost_index"`
		} `json:"minus_6000ft"`
		Minus4000Ft struct {
			TimeEnroute    string `json:"time_enroute"`
			TimeDifference string `json:"time_difference"`
			EnrouteBurn    string `json:"enroute_burn"`
			BurnDifference string `json:"burn_difference"`
			RampFuel       string `json:"ramp_fuel"`
			InitialFl      string `json:"initial_fl"`
			InitialTAS     string `json:"initial_tas"`
			InitialMach    string `json:"initial_mach"`
			CostIndex      string `json:"cost_index"`
		} `json:"minus_4000ft"`
		Minus2000Ft struct {
			TimeEnroute    string `json:"time_enroute"`
			TimeDifference string `json:"time_difference"`
			EnrouteBurn    string `json:"enroute_burn"`
			BurnDifference string `json:"burn_difference"`
			RampFuel       string `json:"ramp_fuel"`
			InitialFl      string `json:"initial_fl"`
			InitialTAS     string `json:"initial_tas"`
			InitialMach    string `json:"initial_mach"`
			CostIndex      string `json:"cost_index"`
		} `json:"minus_2000ft"`
		Plus2000Ft  FlexibleField `json:"plus_2000ft"`
		Plus4000Ft  FlexibleField `json:"plus_4000ft"`
		Plus6000Ft  FlexibleField `json:"plus_6000ft"`
		HigherCi    FlexibleField `json:"higher_ci"`
		LowerCi     FlexibleField `json:"lower_ci"`
		ZfwPlus1000 struct {
			TimeEnroute    string `json:"time_enroute"`
			TimeDifference string `json:"time_difference"`
			EnrouteBurn    string `json:"enroute_burn"`
			BurnDifference string `json:"burn_difference"`
			RampFuel       string `json:"ramp_fuel"`
			InitialFl      string `json:"initial_fl"`
			InitialTAS     string `json:"initial_tas"`
			InitialMach    string `json:"initial_mach"`
			CostIndex      string `json:"cost_index"`
		} `json:"zfw_plus_1000"`
		ZfwMinus1000 struct {
			TimeEnroute    string `json:"time_enroute"`
			TimeDifference string `json:"time_difference"`
			EnrouteBurn    string `json:"enroute_burn"`
			BurnDifference string `json:"burn_difference"`
			RampFuel       string `json:"ramp_fuel"`
			InitialFl      string `json:"initial_fl"`
			InitialTAS     string `json:"initial_tas"`
			InitialMach    string `json:"initial_mach"`
			CostIndex      string `json:"cost_index"`
		} `json:"zfw_minus_1000"`
	} `json:"impacts"`
	Crew struct {
		PilotId string `json:"pilot_id"`
		Cpt     string `json:"cpt"`
		Fo      string `json:"fo"`
		Dx      string `json:"dx"`
		Pu      string `json:"pu"`
	} `json:"crew"`
	Weather struct {
		OrigMETAR   string        `json:"orig_metar"`
		OrigTAF     string        `json:"orig_taf"`
		DestMETAR   string        `json:"dest_metar"`
		DestTAF     string        `json:"dest_taf"`
		AltnMETAR   string        `json:"altn_metar"`
		AltnTAF     string        `json:"altn_taf"`
		ToaltnMETAR FlexibleField `json:"toaltn_metar"`
		ToaltnTAF   FlexibleField `json:"toaltn_taf"`
		EualtnMETAR FlexibleField `json:"eualtn_metar"`
		EualtnTAF   FlexibleField `json:"eualtn_taf"`
		EtopsMETAR  FlexibleField `json:"etops_metar"`
		EtopsTAF    FlexibleField `json:"etops_taf"`
	} `json:"weather"`
	Tracks          FlexibleField `json:"tracks"`
	DatabaseUpdates struct {
		METARTAF string `json:"metar_taf"`
		Winds    string `json:"winds"`
		Sigwx    string `json:"sigwx"`
		Sigmet   string `json:"sigmet"`
		Notams   string `json:"notams"`
		Tracks   string `json:"tracks"`
	} `json:"database_updates"`
	Links struct {
		Skyvector string `json:"skyvector"`
	} `json:"links"`
	VatsimPrefile    string `json:"vatsim_prefile"`
	IvaoPrefile      string `json:"ivao_prefile"`
	PilotedgePrefile string `json:"pilotedge_prefile"`
	PosconPrefile    string `json:"poscon_prefile"`
	MapData          string `json:"map_data"`
	ApiParams        struct {
		Airline        FlexibleField `json:"airline"`
		Fltnum         string        `json:"fltnum"`
		Type           string        `json:"type"`
		Orig           string        `json:"orig"`
		Dest           string        `json:"dest"`
		Date           string        `json:"date"`
		Dephour        string        `json:"dephour"`
		Depmin         string        `json:"depmin"`
		Route          string        `json:"route"`
		Stehour        string        `json:"stehour"`
		Stemin         string        `json:"stemin"`
		Reg            string        `json:"reg"`
		Fin            string        `json:"fin"`
		Selcal         string        `json:"selcal"`
		Pax            string        `json:"pax"`
		Altn           string        `json:"altn"`
		Fl             FlexibleField `json:"fl"`
		Cpt            string        `json:"cpt"`
		Pid            string        `json:"pid"`
		Fuelfactor     string        `json:"fuelfactor"`
		Manualpayload  string        `json:"manualpayload"`
		Manualzfw      string        `json:"manualzfw"`
		Taxifuel       string        `json:"taxifuel"`
		Minfob         string        `json:"minfob"`
		MinfobUnits    string        `json:"minfob_units"`
		Minfod         string        `json:"minfod"`
		MinfodUnits    string        `json:"minfod_units"`
		Melfuel        string        `json:"melfuel"`
		MelfuelUnits   string        `json:"melfuel_units"`
		Atcfuel        string        `json:"atcfuel"`
		AtcfuelUnits   string        `json:"atcfuel_units"`
		Wxxfuel        string        `json:"wxxfuel"`
		WxxfuelUnits   string        `json:"wxxfuel_units"`
		Addedfuel      string        `json:"addedfuel"`
		AddedfuelUnits string        `json:"addedfuel_units"`
		AddedfuelLabel string        `json:"addedfuel_label"`
		Tankering      string        `json:"tankering"`
		TankeringUnits string        `json:"tankering_units"`
		Flightrules    string        `json:"flightrules"`
		Flighttype     string        `json:"flighttype"`
		Contpct        string        `json:"contpct"`
		Resvrule       string        `json:"resvrule"`
		Taxiout        string        `json:"taxiout"`
		Taxiin         string        `json:"taxiin"`
		Cargo          string        `json:"cargo"`
		Origrwy        string        `json:"origrwy"`
		Destrwy        string        `json:"destrwy"`
		Climb          string        `json:"climb"`
		Descent        string        `json:"descent"`
		Cruisemode     string        `json:"cruisemode"`
		Cruisesub      string        `json:"cruisesub"`
		Planformat     string        `json:"planformat"`
		Pounds         string        `json:"pounds"`
		Navlog         string        `json:"navlog"`
		Etops          string        `json:"etops"`
		Stepclimbs     string        `json:"stepclimbs"`
		Tlr            string        `json:"tlr"`
		NotamsOpt      string        `json:"notams_opt"`
		Firnot         string        `json:"firnot"`
		Maps           string        `json:"maps"`
		Turntoflt      FlexibleField `json:"turntoflt"`
		Turntoapt      FlexibleField `json:"turntoapt"`
		Turntotime     FlexibleField `json:"turntotime"`
		Turnfrflt      FlexibleField `json:"turnfrflt"`
		Turnfrapt      FlexibleField `json:"turnfrapt"`
		Turnfrtime     FlexibleField `json:"turnfrtime"`
		Fuelstats      FlexibleField `json:"fuelstats"`
		Contlabel      FlexibleField `json:"contlabel"`
		StaticId       FlexibleField `json:"static_id"`
		Acdata         FlexibleField `json:"acdata"`
		AcdataParsed   string        `json:"acdata_parsed"`
	} `json:"api_params"`
}
