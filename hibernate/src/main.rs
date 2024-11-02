use clap::ArgAction::HelpLong;
use std::error::Error;
use clap::Parser;
use windows::Win32::Foundation::BOOLEAN;
use windows::Win32::System::Power::SetSuspendState;

/// 通过关闭电源来暂停系统。 根据 `hibernate` 参数，系统进入暂停 (睡眠) 状态或休眠 (S4) 。
///
/// # 参数
///
/// * `hibernate` - 如果此参数为 `true`，则系统休眠。 如果参数为 `false`，系统会挂起。
/// * `force` - 此参数不起作用。
/// * `wakeup_events_disabled` - 如果此参数为 `true`，系统将禁用所有唤醒事件。 如果参数为 `false`，则保持启用任何系统唤醒事件。
fn set_suspend_state(
    hibernate: bool,
    force: bool,
    wakeup_events_disabled: bool,
) -> Result<(), Box<dyn Error>> {
    unsafe {
        SetSuspendState::<BOOLEAN, BOOLEAN, BOOLEAN>(
            hibernate.into(),
            force.into(),
            wakeup_events_disabled.into(),
        )
    }.ok().map_err(|err| err.into())
}

#[derive(Parser, Debug)]
#[command(disable_help_flag = true)]
struct Args {
    /// Make the system hibernate or suspend
    #[arg(short, long, default_value_t = true)]
    hibernate: bool,

    /// Disable all wake events
    #[arg(short, long, default_value_t = false)]
    wakeup_events_disabled: bool,

    /// Print help
    #[arg(long, action = HelpLong)]
    help: Option<bool>,
}

fn main() {
    let args = Args::parse();

    set_suspend_state(args.hibernate, false, args.wakeup_events_disabled).unwrap();
}
