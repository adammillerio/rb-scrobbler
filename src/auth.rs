use std::fs;
extern crate dirs;
use dirs::*;
extern crate rustfm_scrobble;
use rustfm_scrobble::Scrobbler;

pub fn initial_authentication(scrob: &mut Scrobbler, user: &str, pass: &str) {
    let auth_res = scrob.authenticate_with_password(user, pass);
    if auth_res.is_ok() {
        println!("Authenticated!");

        let mut config_directory = config_dir()
            .expect("Couldn't get config directory!");

        config_directory.push("rb-scrobbler");

        fs::create_dir(config_directory)
            .expect("Couldn't create config directory!");

        let mut key_directory = config_dir()
            .expect("Couldn't get config directory!");

        key_directory.push("rb-scrobbler");
        key_directory.push("session");
        key_directory.set_extension("key");

        let session_key = scrob.session_key()
            .expect("Couldn't get session key!");

        fs::write(key_directory, session_key)
            .expect("Couldn't save session key to directory!");

    } else {
        panic!("Authentication failed!");
    }
}
