package gateway

import (
	"fmt"
	"strconv"
	"vehicles/packages/domain/models"

	"github.com/lib/pq"
)

// SelectCars получает из реляционной БД под управлением PostgreSQL информацию об автомобилях
// Входной параметр: sln - запрос пользователя
func (slr *selectionRepository) SelectCars(sln models.Selection) ([]models.Car, error) {
	query := `SELECT makes.make, models.model, generations.generation, steering_wheel_positions.position, power_steering_types.power_steering, 
	body_types.body, specifications.length, specifications.width, specifications.height, specifications.ground_clearance, 
	specifications.drag_coefficient, specifications.front_track_width, specifications.back_track_width, specifications.wheelbase, 
	specifications.crash_test_estimate, specifications.year, engines.fuel_used, engines.engine_type, engines.capacity, engines.power, 
	engines.max_torque, gearboxes.gearbox, drive_types.drive, suspensions.front_stabilizer, suspensions.back_stabilizer, 
	suspensions.front_suspension, suspensions.back_suspension, tires.back_tires_width, tires.front_tires_width, tires.front_tires_aspect_ratio, 
	tires.back_tires_aspect_ratio, tires.front_tires_rim_diameter, tires.back_tires_rim_diameter, brakes.front_brakes, brakes.back_brakes, 
	brakes.parking_brake, safety_and_motion_control_systems.abs_system, safety_and_motion_control_systems.esp_system, 
	safety_and_motion_control_systems.ebd_system, safety_and_motion_control_systems.bas_system,
	safety_and_motion_control_systems.tcs_system, safety_and_motion_control_systems.front_parking_sensor, 
	safety_and_motion_control_systems.back_parking_sensor, safety_and_motion_control_systems.rear_view_camera, 
	safety_and_motion_control_systems.cruise_control, colors.color, lights.headlights, lights.led_running_lights, lights.led_tail_lights, 
	lights.light_sensor, lights.front_fog_lights, lights.back_fog_lights, interior_design.upholstery, cabin_microclimate.air_conditioner, 
	cabin_microclimate.climate_control, electric_options.electric_front_side_windows_lifts, electric_options.electric_back_side_windows_lifts, 
	electric_options.electric_heating_of_front_seats, electric_options.electric_heating_of_back_seats, 
	electric_options.electric_heating_of_steering_wheel, electric_options.electric_heating_of_windshield, 
	electric_options.electric_heating_of_rear_window, electric_options.electric_heating_of_side_mirrors, 
	electric_options.electric_drive_of_driver_seat, electric_options.electric_drive_of_front_seats,
	electric_options.electric_drive_of_side_mirrors, electric_options.electric_trunk_opener, electric_options.rain_sensor, 
	airbags.driver_airbag, airbags.front_passenger_airbag, airbags.side_airbags, airbags.curtain_airbags, multimedia_systems.on_board_computer, 
	multimedia_systems.mp3_support, multimedia_systems.hands_free_support, trim_levels.trim_level, trim_levels.acceleration_0_to_100, 
	trim_levels.max_speed, trim_levels.city_fuel_consumption, trim_levels.highway_fuel_consumption, trim_levels.mixed_fuel_consumption, 
	trim_levels.number_of_seats, trim_levels.trunk_volume, trim_levels.mass, trim_levels.car_alarm, offerings.price, 
	offerings.kilometerage, offerings.photo_urls
		FROM makes
		INNER JOIN countries ON makes.country_id = countries.id
		INNER JOIN models ON makes.id = models.make_id
		INNER JOIN generations ON models.id = generations.model_id
		INNER JOIN specifications ON generations.id = specifications.generation_id
		INNER JOIN trim_levels ON specifications.id = trim_levels.specification_id
		INNER JOIN engines ON trim_levels.engine_id = engines.id
		INNER JOIN gearboxes ON trim_levels.gearbox_id = gearboxes.id
		INNER JOIN drive_types ON trim_levels.drive_type_id = drive_types.id
		INNER JOIN suspensions ON specifications.id = suspensions.id
		INNER JOIN tires ON specifications.id = tires.id
		INNER JOIN brakes ON specifications.id = brakes.id
		INNER JOIN safety_and_motion_control_systems ON trim_levels.safety_and_motion_control_systems_id = safety_and_motion_control_systems.id
		INNER JOIN colors ON trim_levels.color_id = colors.id
		INNER JOIN lights ON trim_levels.lights_id = lights.id
		INNER JOIN cabin_microclimate ON trim_levels.cabin_microclimate_id = cabin_microclimate.id
		INNER JOIN electric_options ON trim_levels.electric_options_id = electric_options.id
		INNER JOIN multimedia_systems ON trim_levels.multimedia_systems_id = multimedia_systems.id
		INNER JOIN offerings ON trim_levels.id = offerings.trim_level_id
		LEFT JOIN steering_wheel_positions ON specifications.steering_wheel_position_id = steering_wheel_positions.id
		LEFT JOIN power_steering_types ON specifications.power_steering_type_id = power_steering_types.id
		LEFT JOIN body_types ON specifications.body_type_id = body_types.id
		LEFT JOIN interior_design ON trim_levels.interior_design_id = interior_design.id
		LEFT JOIN airbags ON trim_levels.airbags_id = airbags.id`

	// добавление условия по странам и цене
	whereClause := ""
	args := make([]interface{}, 0)

	if sln.MinPrice != "" || sln.MaxPrice != "" {
		where := func() string {
			if whereClause == "" {
				return "WHERE "
			}
			return " AND "
		}()
		whereClause = fmt.Sprintf("%s%s", whereClause, where)

		switch {
		case sln.MinPrice != "" && sln.MaxPrice != "":
			whereClause = fmt.Sprintf("%sofferings.price BETWEEN $1 AND $2", whereClause)
			args = append(args, sln.MinPrice, sln.MaxPrice)
		case sln.MinPrice != "":
			whereClause = fmt.Sprintf("%sofferings.price >= $1", whereClause)
			args = append(args, sln.MinPrice)
		case sln.MaxPrice != "":
			whereClause = fmt.Sprintf("%sofferings.price <= $1", whereClause)
			args = append(args, sln.MaxPrice)
		}
	}

	for i, man := range sln.Manufacturers {
		if man == "Другие" {
			sln.Manufacturers[i] = "Чехия"
		}
	}

	if len(sln.Manufacturers) != 0 {
		if whereClause == "" {
			whereClause = fmt.Sprintf("%sWHERE ", whereClause)
		} else {
			whereClause = fmt.Sprintf("%s AND ", whereClause)
		}
		whereClause = fmt.Sprintf("%scountries.country IN (", whereClause)
		for i, m := range sln.Manufacturers {
			args = append(args, m)
			whereClause = fmt.Sprintf("%s$%d", whereClause, len(args))
			if i < len(sln.Manufacturers)-1 {
				whereClause = fmt.Sprintf("%s, ", whereClause)
			}
		}
		whereClause = fmt.Sprintf("%s)", whereClause)
	}

	if whereClause != "" {
		query = fmt.Sprintf("%s %s", query, whereClause)
	}

	rows, err := slr.vehiclesDB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error from `Query` method, package `sql`: %#v", err)
	}
	defer rows.Close()

	cars := []models.Car{}
	index := 0
	for rows.Next() {
		car := models.NewCar()
		var make, model string
		err := rows.Scan(&make, &model, &car.Generation, &car.Specs.SteeringWheel.SteeringWheelPosition, &car.Specs.SteeringWheel.PowerSteering,
			&car.Specs.Body, &car.Specs.Length, &car.Specs.Width, &car.Specs.Height, &car.Specs.GroundClearance, &car.Specs.DragCoefficient,
			&car.Specs.FrontTrackWidth, &car.Specs.BackTrackWidth, &car.Specs.Wheelbase, &car.Specs.CrashTestEstimate, &car.Offering.Year,
			&car.Specs.Engine.FuelUsed, &car.Specs.Engine.EngineType, &car.Specs.Engine.Capacity, &car.Specs.Engine.MaxPower,
			&car.Specs.Engine.MaxTorque, &car.Specs.Gearbox, &car.Specs.Drive, &car.Specs.Suspension.FrontStabilizer,
			&car.Specs.Suspension.BackStabilizer, &car.Specs.Suspension.FrontSuspension, &car.Specs.Suspension.BackSuspension,
			&car.Specs.Tires.BackTiresWidth, &car.Specs.Tires.FrontTiresWidth, &car.Specs.Tires.FrontTiresAspectRatio,
			&car.Specs.Tires.BackTiresAspectRatio, &car.Specs.Tires.FrontTiresRimDiameter, &car.Specs.Tires.BackTiresRimDiameter,
			&car.Specs.Brakes.FrontBrakes, &car.Specs.Brakes.BackBrakes, &car.Specs.Brakes.ParkingBrake,
			&car.Features.SafetyAndMotionControlSystem.ABS, &car.Features.SafetyAndMotionControlSystem.ESP,
			&car.Features.SafetyAndMotionControlSystem.EBD, &car.Features.SafetyAndMotionControlSystem.BAS,
			&car.Features.SafetyAndMotionControlSystem.TCS, &car.Features.SafetyAndMotionControlSystem.FrontParkingSensor,
			&car.Features.SafetyAndMotionControlSystem.BackParkingSensor, &car.Features.SafetyAndMotionControlSystem.RearViewCamera,
			&car.Features.SafetyAndMotionControlSystem.CruiseControl, &car.Features.Color, &car.Features.Lights.Headlights,
			&car.Features.Lights.LEDRunningLights, &car.Features.Lights.LEDTailLights, &car.Features.Lights.LightSensor,
			&car.Features.Lights.FrontFogLights, &car.Features.Lights.BackFogLights, &car.Features.Interior.Upholstery,
			&car.Features.CabinMicroclimate.AirConditioner, &car.Features.CabinMicroclimate.ClimateControl,
			&car.Features.ElectricOptions.ElectricFrontSideWindowsLifts, &car.Features.ElectricOptions.ElectricBackSideWindowsLifts,
			&car.Features.ElectricOptions.ElectricHeatingOfFrontSeats, &car.Features.ElectricOptions.ElectricHeatingOfBackSeats,
			&car.Features.ElectricOptions.ElectricHeatingOfSteeringWheel, &car.Features.ElectricOptions.ElectricHeatingOfWindshield,
			&car.Features.ElectricOptions.ElectricHeatingOfRearWindow, &car.Features.ElectricOptions.ElectricHeatingOfSideMirrors,
			&car.Features.ElectricOptions.ElectricDriveOfDriverSeat, &car.Features.ElectricOptions.ElectricDriveOfFrontSeats,
			&car.Features.ElectricOptions.ElectricDriveOfSideMirrors, &car.Features.ElectricOptions.ElectricTrunkOpener,
			&car.Features.ElectricOptions.RainSensor, &car.Features.Airbags.DriverAirbag, &car.Features.Airbags.FrontPassengerAirbag,
			&car.Features.Airbags.SideAirbags, &car.Features.Airbags.CurtainAirbags, &car.Features.MultimediaSystems.OnBoardComputer,
			&car.Features.MultimediaSystems.MP3Support, &car.Features.MultimediaSystems.HandsFreeSupport, &car.TrimLevel,
			&car.Specs.Acceleration0To100, &car.Specs.MaxSpeed, &car.Specs.CityFuelConsumption, &car.Specs.HighwayFuelConsumption,
			&car.Specs.MixedFuelConsumption, &car.Specs.NumberOfSeats, &car.Specs.TrunkVolume, &car.Specs.Mass, &car.Features.CarAlarm,
			&car.Offering.Price, &car.Offering.Kilometerage,
			pq.Array(&car.Offering.PhotoURLs))

		car.ID = index
		car.FullName = fmt.Sprintf("%s %s, %s", make, model, strconv.Itoa(car.Offering.Year))
		car.Offering.Price = formatPrice(car.Offering.Price)
		cars = append(cars, car)
		index++
		if err != nil {
			return nil, fmt.Errorf("error from `Scan` method, package `sql`: %#v", err)
		}
	}
	return cars, nil
}

// formatPrice приводит цену к формату
// Входной параметр: price - цена
func formatPrice(price string) string {
	var formattedPrice string
	for i := len(price) - 3; i > 0; i -= 3 {
		formattedPrice = fmt.Sprintf("%s %s", price[:i], price[i:])
	}

	formattedPrice = fmt.Sprintf("%s ₽", formattedPrice)
	return formattedPrice
}
