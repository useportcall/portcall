<#macro emailLayout>
<!DOCTYPE html>
<html lang="${locale.language}" dir="${(ltr)?then('ltr','rtl')}">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>Portcall</title>
    <style>
        body {
            margin: 0;
            padding: 0;
            background-color: #f4f4f5;
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
        }
        .email-wrapper {
            width: 100%;
            background-color: #f4f4f5;
            padding: 40px 0;
        }
        .email-container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #ffffff;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
        }
        .email-header {
            background-color: #18181b;
            padding: 24px 40px;
            text-align: center;
        }
        .logo {
            font-size: 24px;
            font-weight: 700;
            color: #ffffff;
            text-decoration: none;
            letter-spacing: -0.5px;
        }
        .logo-accent {
            color: #6366f1;
        }
        .email-body {
            padding: 40px;
        }
        .email-body h1 {
            color: #18181b;
            font-size: 24px;
            font-weight: 600;
            margin: 0 0 16px 0;
        }
        .email-body p {
            color: #52525b;
            font-size: 16px;
            line-height: 1.6;
            margin: 0 0 16px 0;
        }
        .email-body a {
            color: #6366f1;
            text-decoration: none;
        }
        .email-body a:hover {
            text-decoration: underline;
        }
        .button {
            display: inline-block;
            background-color: #18181b;
            color: #ffffff !important;
            padding: 14px 32px;
            border-radius: 12px;
            font-size: 16px;
            font-weight: 600;
            text-decoration: none;
            margin: 16px 0;
        }
        .button:hover {
            background-color: #27272a;
            text-decoration: none;
        }
        .email-footer {
            background-color: #fafafa;
            padding: 24px 40px;
            text-align: center;
            border-top: 1px solid #e4e4e7;
        }
        .email-footer p {
            color: #a1a1aa;
            font-size: 14px;
            margin: 0;
        }
        .email-footer a {
            color: #71717a;
            text-decoration: none;
        }
    </style>
</head>
<body>
    <div class="email-wrapper">
        <div class="email-container">
            <div class="email-header">
                <span class="logo">Port<span class="logo-accent">call</span></span>
            </div>
            <div class="email-body">
                <#nested>
            </div>
            <div class="email-footer">
                <p>&copy; ${.now?string('yyyy')} Portcall. All rights reserved.</p>
                <p><a href="https://useportcall.com">useportcall.com</a></p>
            </div>
        </div>
    </div>
</body>
</html>
</#macro>
