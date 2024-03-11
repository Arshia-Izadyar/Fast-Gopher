ALTER TABLE active_devices
ADD CONSTRAINT unique_device_user UNIQUE (device_id, user_id);
