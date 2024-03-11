ALTER TABLE active_devices
DROP CONSTRAINT unique_device_user UNIQUE (device_id, user_id);
