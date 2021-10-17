//Prompt is the JavaScript module for all alerts, notifications and custom popup dialogs
function Prompt() {

    let toast = function (c) {
        const {
            title = '',
            icon = 'success',
            position = 'top-end'
        } = c
        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            showConfirmButton: false,
            timer: 3000,
            icon: icon,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = function (c) {
        const {
            title = '',
            text = '',
            footer = ''
        } = c
        Swal.fire({
            icon: 'success',
            title: title,
            text: text,
            footer: footer
        })
    }

    let error = function (c) {
        const {
            title = '',
            text = '',
            footer = ''
        } = c
        Swal.fire({
            icon: 'error',
            title: title,
            text: text,
            footer: footer
        })
    }

    async function custom(c) {
        const {
            icon = '',
            title = '',
            html = '',
            showConfirmButton = true
        } = c
        const {value: result} = await Swal.fire({
            icon: icon,
            title: title,
            html: html,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            },
            preConfirm: () => {
                return [
                    document.getElementById('start_date').value,
                    document.getElementById('end_date').value
                ]
            }
        })
        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result !== "") {
                    if (c.callback !== undefined) {
                        c.callback(result);
                    } else {
                        c.callback(false);
                    }
                }
            } else {
                c.callback(false);
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom
    }
}