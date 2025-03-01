window.onload = () => {
    let isWafEnabled = false
    const wafCheckbox = document.getElementById("waf-checkbox")
    const title = document.getElementById("title")

    const onWafEnabled = (checked) => {
        isWafEnabled = checked
        if (isWafEnabled) {
            title.classList.add('text-success')
            title.classList.remove('text-fail')
            title.innerHTML = 'is enabled'
        } else {
            title.classList.add('text-fail')
            title.classList.remove('text-success')
            title.innerHTML = 'is not enabled'
        }
    }

    wafCheckbox.onchange = (e) => onWafEnabled(e.target.checked)
    onWafEnabled(false)

    const lastParameterEl = document.getElementById("last-parameter")
    lastParameterEl.innerHTML = "Last parameter: -"
    const updateLastParameter = (param) => {
        lastParameterEl.innerHTML = "Last parameter: " + param
    }

    const sqliParamInput = document.getElementById("sqli-param")
    const sqliSubmit = document.getElementById("sqli-submit")

    sqliSubmit.onclick = () => {
        const sqliParam = sqliParamInput.value
        const method = 'POST'
        const uri = '/protected'

        fetch('/protected', {
            method: method,
            body: JSON.stringify({
                username: sqliParam,
                enable_waf: isWafEnabled,
            }),
            headers: {
                "Content-Type":"application/json",
            }
        }).then(r => {
            updateLastParameter(sqliParam)
            insertRequest({
                method: method,
                uri: uri,
                param: sqliParam,
                status: r.status,
            })
        })
    }

    const xssParamInput = document.getElementById("xss-param")
    const xssSubmit = document.getElementById("xss-submit")

    xssSubmit.onclick = () => {
        const xssParam = xssParamInput.value
        const method = 'GET'
        const uri = `/protected?enable_waf=${isWafEnabled}&username=${xssParam}`

        fetch(uri, {
            method: method
        }).then(r => {
            updateLastParameter(xssParam)
            insertRequest({
                method: method,
                uri: uri,
                param: xssParam,
                status: r.status,
            })
        })
    }

    const requestHistory = []
    const requestTable = document.getElementById("request-table")
    const insertRequest = (params) => {
        const row = requestTable.insertRow()

        requestHistory.push({
            method: params.method,
            uri: params.uri,
            param: params.param,
            status: params.status,
        })

        const method = row.insertCell(0)
        method.innerHTML = params.method
        method.classList.add('attacks-table__cell')
        const uri = row.insertCell(1)
        uri.innerHTML = params.uri
        uri.classList.add('attacks-table__cell')
        const param = row.insertCell(2)
        param.innerHTML = params.param
        param.classList.add('attacks-table__cell')
        const status = row.insertCell(3)
        status.innerHTML = params.status
        status.classList.add('attacks-table__cell')
        if (params.status === 200) {
            status.classList.add('bg-success')
        } else {
            status.classList.add('bg-fail')
        }
    }
}