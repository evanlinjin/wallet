"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
var core_1 = require("@angular/core");
var common_1 = require("@angular/common");
var forms_1 = require("@angular/forms");
var core_2 = require("@ngx-translate/core");
var ng_bootstrap_1 = require("@ng-bootstrap/ng-bootstrap");
var register_routing_module_1 = require("./register-routing.module");
var register_component_1 = require("./register.component");
var RegisterModule = /** @class */ (function () {
    function RegisterModule() {
    }
    RegisterModule = __decorate([
        core_1.NgModule({
            imports: [
                common_1.CommonModule,
                forms_1.ReactiveFormsModule,
                core_2.TranslateModule,
                ng_bootstrap_1.NgbModule,
                register_routing_module_1.RegisterRoutingModule
            ],
            declarations: [
                register_component_1.RegisterComponent
            ]
        })
    ], RegisterModule);
    return RegisterModule;
}());
exports.RegisterModule = RegisterModule;
//# sourceMappingURL=register.module.js.map